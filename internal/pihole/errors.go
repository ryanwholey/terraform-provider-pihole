package pihole

import "errors"

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

var (
	// ErrLoginFailed is returned when a login attempt fails
	ErrLoginFailed = errors.New("login failed")
	// ErrClientValidationFailed
	ErrClientValidationFailed = errors.New("client validation failed")
	// ErrNotImplementedTokenClient is returned when a particular Pi-hole resource cannot be managed due to missing client configuration
	ErrNotImplementedTokenClient = errors.New("resource is not implemented for the API token client")
)
