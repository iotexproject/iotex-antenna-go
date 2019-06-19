package errcodes

import "errors"

// Code is a error code
type Code int

// Codes
const (
	InvalidParam Code = iota + 1
	RPCError
	BadResponse
	InternalError
)

// ErrorWithCode is an error with an associated code.
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

// NewError takes an error and a code return a error associated with the code.
func NewError(err error, c Code) error {
	return &ewc{
		err:  err,
		code: c,
	}
}

// New takes an error message string and a code return a error associated with the code.
func New(msg string, c Code) error {
	return &ewc{
		err:  errors.New(msg),
		code: c,
	}
}
