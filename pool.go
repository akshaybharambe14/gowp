package gowp

import (
	"context"
	"errors"
	"sync"
)

type (
	// workerPool represents a pool of workers that limits concurency as per the provided worker count.
	// This is not exported as we want users to use New() to create a worker pool and limit the scope of initialized pool in the same function where it is initialized.
	workerPool struct {
		errOnce *sync.Once    // errOnce ensures that the error is assigned only once.
		errMtx  *sync.RWMutex // errMtx protects the access to err.
		err     error         // err is the first error that occurred in the work.

		wg *sync.WaitGroup

		closeOnce *sync.Once    // closeOnce ensures that close on in and closed is called only once.
		in        chan work     // in works as a queue of work that workers listen on.
		closed    chan struct{} // closed helps to determine if the pool is closed.
	}

	// work is a unit of work that is submitted to the pool by users.
	work func() error
)

func New(ctx context.Context, buffer int, workers int, closeOnErr bool) (*workerPool, error) {
	if buffer <= 0 {
		return nil, ErrInvalidBuffer
	}

	if workers <= 0 {
		return nil, ErrInvalidWorkerCnt
	}

	if ctx == nil {
		ctx = context.TODO()
	}

	return newWorkerPool(ctx, buffer, workers, closeOnErr), nil
}

func newWorkerPool(ctx context.Context, buffer int, workers int, closeOnErr bool) *workerPool {
	wp := &workerPool{
		in:        make(chan work, buffer),
		wg:        &sync.WaitGroup{},
		closeOnce: &sync.Once{},
		errOnce:   &sync.Once{},
		errMtx:    &sync.RWMutex{},
		closed:    make(chan struct{}),
	}

	for i := 0; i < workers; i++ {
		wp.wg.Add(1)

		go func() {
			wp.setError(worker(ctx, wp.in, closeOnErr))
			wp.wg.Done()
		}()
	}

	// this block ensures every goroutine exits before (*workerPool).Wait() returns and we don't accept further requests.
	{
		wp.wg.Add(1)
		go func() {
			defer wp.wg.Done()
			select {
			case <-ctx.Done():
				// user cancelled the context
				// hence ensure that user does not submit any more work to the pool
				wp.Close()

			case <-wp.closed: // the pool is closed
			}
		}()
	}

	return wp
}

// Submit submits work to the pool. It returns false if
// 	1. the pool is closed
// 	2. the passed closure is nil
// 	3. the pool has insufficient buffer
//
// If it returns true, the work is submitted to the pool and will get eventually executed.
func (wp *workerPool) Submit(w work) bool {
	return w != nil && wp.submit(w) == nil
}

// submit tries to submit the work to the pool. It returns nil if the work is successfully submitted.
func (wp *workerPool) submit(w work) (err error) {
	if wp.isClosed() {
		return ErrPoolClosed
	}

	defer func() {
		// at this point, the pool was not closed till it passed the above check, and pool got closed just before send operation. This block ensures that this function does not panic because of send on closed channel.
		if recover() != nil {
			err = ErrInvalidSend
		}
	}()

	select {
	case wp.in <- w:
		// we have sufficient buffer to push work
		return nil

	default:
		// insufficient buffer, return error
		return ErrNoBuffer
	}
}

// Close closes the worker pool. It does not wait for the work to be finished.
//
// It is safe to call this function concurrently.
func (wp *workerPool) Close() {
	wp.closeOnce.Do(func() {
		close(wp.in)
		close(wp.closed)
	})
}

// isClosed returns true if the pool is closed. This means (*workerpool).Close() was called before.
func (wp *workerPool) isClosed() bool {
	select {
	case <-wp.closed:
		return true

	default:
		return false
	}
}

// Wait waits for all the work to be finished.
// It returns the first error if opted for close on error occurred, if any.
func (wp *workerPool) Wait() error {
	/*
		wp.errMtx.RLock()
		defer wp.errMtx.RUnlock()

		wp.wg.Wait()

		return wp.err

		FUN FACT: above code will panic as mutex is acquired and we are waiting to finish the work.
		If any worker won't be able to submit an error because of mutex.DEADLOCK!

	*/

	wp.wg.Wait()

	wp.errMtx.RLock()
	defer wp.errMtx.RUnlock()

	return wp.err
}

func (wp *workerPool) setError(err error) {
	wp.errOnce.Do(func() {
		if isWorkError(err) {
			wp.errMtx.Lock()
			defer wp.errMtx.Unlock()

			wp.err = err
			wp.Close()
		}
	})
}

// worker processes work sent on in channel. When it exists, it is guaranteed that last processed work is finished.
// It must be run in a separate goroutine.
//
// The worker returns when the context is canceled, exceeds deadline OR the pool is closed.
//
// The returned error helps in testing this function.
func worker(ctx context.Context, in chan work, closeOnErr bool) error {
	for {
		select {
		case <-ctx.Done():
			// user cancelled the context
			return ctx.Err()

		case w, ok := <-in:
			if !ok {
				// the channel is closed by calling Close()
				return ErrPoolClosed
			}

			if err := w(); closeOnErr && err != nil {
				return err
			}
		}
	}
}

// isWorkError returns true if the err is a error returned by user defined work function and is not ErrPoolClosed
func isWorkError(err error) bool {
	return !errors.Is(err, ErrPoolClosed)
}
