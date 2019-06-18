package errcodes

import "errors"

type Code int

const (
	InvalidParam Code = iota + 1
	RPCError
	BadResponse
	InternalError
)

type ErrorWithCode interface {
	Error() string
	Code() Code
	Cause() error
}

type ewc struct {
	err  error
	code Code
}

func (e *ewc) Error() string { return e.err.Error() }
func (e *ewc) Code() Code    { return e.code }
func (e *ewc) Cause() error  { return e.err }

func NewError(err error, c Code) error {
	return &ewc{
		err:  err,
		code: c,
	}
}

func New(msg string, c Code) error {
	return &ewc{
		err:  errors.New(msg),
		code: c,
	}
}
