package partition

import "fmt"

// ValidationError represents an error that occured
// after validating a struct's attributes
type ValidationError struct {
	Err error
}

// Returns a formatted error message including the underlying
// error message
func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error: error=%v", e.Err)
}

// Returns the error
func (e *ValidationError) Unwrap() error {
	return e.Err
}

// SetupPartitionsError represents an error that occured
// after trying to create partitions
//
// It wraps the underlying error for better clarity
type SetupPartitionsError struct {
	Err error
}

// Returns a formatted error message including the underlying
// error message
func (e *SetupPartitionsError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the underlying error for error chaining
func (e *SetupPartitionsError) Unwrap() error {
	return e.Err
}

// PartitionTableCompatibilityError represents an error that occured
// after checking the compatibility of a drive's partition table
// with the installer: only GPT is supported
//
// It wraps the underlying error for better clarity
type PartitionTableCompatibilityError struct {
	Err error
}

// Returns a formatted error message including the underlying
// error message
func (e *PartitionTableCompatibilityError) Error() string {
	return e.Err.Error()
}

// Unwrap returns the underlying error for error chaining
func (e *PartitionTableCompatibilityError) Unwrap() error {
	return e.Err
}
