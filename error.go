package gowp

type Error string

// processing errors
const (
	ErrPoolClosed  = Error("pool is closed")
	ErrNoBuffer    = Error("insufficient buffer")
	ErrInvalidSend = Error("work sent on closed pool")
	ErrNilTask     = Error("task is nil")
)

// validation errors
const (
	ErrInvalidBuffer    = Error("buffer value should be greater than zero")
	ErrInvalidWorkerCnt = Error("worker count should be greater than zero")
	ErrNilContext       = Error("context is nil")
)

// interface guard to ensure Error implements error interface
var _ error = Error("")

func (e Error) Error() string {
	return string(e)
}
