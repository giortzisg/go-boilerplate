package error

type StatusError struct {
	error
	code int
}

func NewStatusError(err error, code int) *StatusError {
	return &StatusError{
		error: err,
		code:  code,
	}
}

func (e *StatusError) Unwrap() error {
	return e.error
}

func (e *StatusError) HTTPStatus() int {
	return e.code
}
