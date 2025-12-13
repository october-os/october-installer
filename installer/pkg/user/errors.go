package user

import "fmt"

// UserInstantiationError represents an error that occured during the
// instantiation of a user inside the constructor.
//
// It wraps an error for better debugging.
type UserInstantiationError struct {
	err error
}

// Error returns a formatted string containing the error message.
func (e UserInstantiationError) Error() string {
	return fmt.Sprintf("failed to instantiate user: error=%v", e.err.Error())
}

// Unwrap unwraps the error inside UserIntantiationError
func (e UserInstantiationError) Unwrap() error {
	return e.err
}
