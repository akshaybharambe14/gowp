package gowp

type Error string

// processing errors
const (
	ErrPoolClosed  = Error("pool is closed")
	ErrNoBuffer    = Error("insufficient buffer")
	ErrInvalidSend = Error("work sent on closed pool")
)

// validation errors
const (
	ErrInvalidBuffer    = Error("buffer value should be more than zero")
	ErrInvalidWorkerCnt = Error("worker count should be more than zero")
)

var _ error = Error("")

func (e Error) Error() string {
	return string(e)
}
