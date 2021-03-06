package gowp

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
)

// closed represents the closed state of the pool.
const closed uint32 = 1

type (
	// Pool represents a pool of workers that limits concurency as per the provided worker count.
	//
	// Zero value is not usable. Use New() to create a new Pool.
	Pool struct {
		wg sync.WaitGroup

		err          error         // the first error that occurred in the execution.
		errs         chan error    // workers report errors through this channel.
		quit         chan struct{} // quit signal to close the pool. This will be closed on error or after successful execution.
		exitFromErrG chan struct{} // exit signal to close the error handling goroutine, in case if not closed already.
		in           chan Task     // works as a queue of work that workers listen to.
		closeOnce    sync.Once     // ensures that we perform exit formalities only once.
		closed       uint32        // set to closed(1) when the pool is closed. Should be manipulated by sync/atomic.

		// Initially, it was thought that not to export this type
		// as we want to force users to use New() to create a new pool
		// and limit the scope of initialized pool to the same function
		// where it is initialized. This turns out to be a bad design.
		// see https://github.com/golang/go/issues/2273
	}

	// Task is a unit of work that is submitted to the pool by consumers.
	Task func() error
)

func New(numTasks int, opts ...Option) (*Pool, error) {
	if numTasks <= 0 {
		return nil, fmt.Errorf("gowp.New(): %w", ErrInvalidBuffer)
	}

	cfg := config{
		ctx:        context.TODO(),
		numWorkers: runtime.NumCPU(),
	}

	for _, opt := range opts {
		opt(&cfg)
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("gowp.New(): %w", err)
	}

	return newPool(cfg.ctx, cfg.numWorkers, numTasks, cfg.exitOnErr), nil
}

func (p *Pool) IsClosed() bool {
	return atomic.LoadUint32(&p.closed) == closed
}

func (p *Pool) Submit(t Task) error {
	if err := p.submit(t); err != nil {
		return fmt.Errorf("gowp.Pool.Submit(): %w", err)
	}

	return nil
}

func (p *Pool) Wait() error {
	p.closeOnce.Do(func() {
		close(p.in)
		atomic.StoreUint32(&p.closed, closed)

		p.wg.Wait() // here, all workers are returned and no worker is writing to p.errs Only error handling go routine will write an error, if any.

		close(p.exitFromErrG) // signal to the error handling go routine to exit (if not initiated by error occurrence OR context cancellation).

		p.err = <-p.errs // wait for the error handling go routine to exit and write an error, if any.
	})

	if p.err != nil {
		return fmt.Errorf("gowp.Pool.Wait(): %w", p.err)
	}

	return nil
}

func newPool(ctx context.Context, numWorkers, numTasks int, exitOnErr bool) *Pool {
	p := &Pool{
		wg:           sync.WaitGroup{},
		in:           make(chan Task, numTasks),
		closeOnce:    sync.Once{},
		errs:         make(chan error, 1),
		quit:         make(chan struct{}, 1),
		exitFromErrG: make(chan struct{}, 1),
	}

	go func() {
		var err error
		select {
		case <-ctx.Done():
			err = ctx.Err()
			close(p.quit)

		case err = <-p.errs:
			if exitOnErr {
				close(p.quit)
			}

		case <-p.exitFromErrG:
			// p.Wait() will be close p.exitFromErrG to signal the exit.
			// this helps to avoid goroutine leak, in case if we don't encounter any errors.
		}

		p.errs <- err
	}()

	for i := 0; i < numWorkers; i++ {
		p.wg.Add(1)
		go func() {
			defer p.wg.Done()
			work(p.in, p.quit, p.errs)
		}()
	}

	return p
}

func (p *Pool) submit(t Task) (err error) {
	if t == nil {
		return ErrNilTask
	}

	if p.IsClosed() {
		return ErrPoolClosed
	}

	defer func() {
		if p := recover(); p != nil {
			err = ErrInvalidSend
		}
	}()

	select {
	case p.in <- t:
		err = nil
	default:
		err = ErrNoBuffer
	}

	return
}

func work(in <-chan Task, quit <-chan struct{}, errs chan<- error) {
	for {
		select {
		case <-quit:
			return
		case t, ok := <-in:
			if !ok {
				return
			}

			if err := t(); err != nil {
				select {
				case errs <- err:
				default:
					// drop the error as p.errs is full, eventually it will receive quit signal
				}
			}
		}
	}
}
