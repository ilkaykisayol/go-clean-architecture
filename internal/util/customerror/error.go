package customerror

type Error struct {
	Loglevel Loglevel
	Err      error
}

// NewErrorWithLoglevel returns a new Error with loglevel.
func New(err error, loglevel Loglevel) *Error {
	return &Error{
		Loglevel: loglevel,
		Err:      err,
	}
}

// Error returns the error message.
func (e *Error) Error() string {
	return e.Err.Error()
}
