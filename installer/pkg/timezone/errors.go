package timezone

import "fmt"

// PipeError represents an error that occurred after
// a failed attempt to pipe STDOUT.
//
// It wraps the original error for better clarity.
type PipeError struct {
	Err error
}

// Returns a formatted error message including the original
// error message.
func (e PipeError) Error() string {
	return fmt.Sprintf("STDOUT pipe creation failed: error=%v", e.Err)
}

// Unwrap returns the underlying error for error chaining.
func (e PipeError) Unwrap() error {
	return e.Err
}

type InvalidTimezoneError struct {
	Err error
}

// Returns a formatted error message including the original
// error message.
func (e InvalidTimezoneError) Error() string {
	return fmt.Sprintf("invalid timezone error: error=%v", e.Err)
}

// Unwrap returns the underlying error for error chaining.
func (e InvalidTimezoneError) Unwrap() error {
	return e.Err
}
