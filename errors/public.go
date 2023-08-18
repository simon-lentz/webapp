package errors

type PublicError interface {
	error
	Public() string
}

// Public wraps the original error with a new error
// with a Public() method, returning a public safe
// error message.
func Public(err error, msg string) error {
	return publicError{err, msg}
}

type publicError struct {
	err error
	msg string
}

func (pe publicError) Error() string {
	return pe.err.Error()
}

func (pe publicError) Public() string {
	return pe.msg
}

func (pe publicError) Unwrap() error {
	return pe.err
}
