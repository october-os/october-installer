package user

import "fmt"

// NewUserError represents an error that occured during the
// instantiation of a user inside the constructor.
//
// It wraps an error for better debugging.
type NewUserError struct {
	err error
}

// Error returns a formatted string containing the error message.
func (e NewUserError) Error() string {
	return fmt.Sprintf("New user error: error=%v", e.err.Error())
}

// Unwrap unwraps the error inside NewUserError
func (e NewUserError) Unwrap() error {
	return e.err
}
