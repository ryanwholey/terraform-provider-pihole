package pihole

type NotFoundError struct {
	err string
}

// New returns a new NotFoundError
func NewNotFoundError(message string) *NotFoundError {
	return &NotFoundError{
		err: message,
	}
}

func (e NotFoundError) Is(target error) bool {
	return target.Error() == e.err
}

func (e *NotFoundError) Error() string {
	return e.err
}
