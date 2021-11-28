package gowp

import "context"

type config struct {
	ctx        context.Context
	numWorkers int
	exitOnErr  bool
}

type Option func(o *config)

// WithContext returns an Option that sets the context for the pool.
// If the context is canceled the pool will be closed.
func WithContext(ctx context.Context) Option {
	return func(o *config) {
		o.ctx = ctx
	}
}

// WithNumWorkers returns an Option that sets the number of workers for the pool.
// If the number of workers is less than or equal to zero, ErrInvalidWorkerCnt will be returned on Poll initialization.
func WithNumWorkers(numWorkers int) Option {
	return func(o *config) {
		o.numWorkers = numWorkers
	}
}

// WithExitOnError returns an Option that sets the exitOnErr for the pool.
// If the exitOnErr is true, the pool will be closed when the first error is received.
func WithExitOnError(exitOnErr bool) Option {
	return func(o *config) {
		o.exitOnErr = exitOnErr
	}
}

func (o *config) validate() error {
	if o.numWorkers <= 0 {
		return ErrInvalidWorkerCnt
	}

	if o.ctx == nil {
		return ErrNilContext
	}

	return nil
}
