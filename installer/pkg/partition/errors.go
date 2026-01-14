package partition

import "fmt"

// PartitionError represents an error that occured
// after trying to create partitions
//
// It wraps the underlying error for better clarity
type PartitionError struct {
	err error
}

// Returns a formatted error message including the underlying
// error message
func (e *PartitionError) Error() string {
	return fmt.Sprintf("Error setting up partitions: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e *PartitionError) Unwrap() error {
	return e.err
}

// PartitionTableCompatibilityError represents an error that occured
// after checking the compatibility of a drive's partition table
// with the installer: only GPT is supported
//
// It wraps the underlying error for better clarity
type PartitionTableCompatibilityError struct {
	err error
}

// Returns a formatted error message including the underlying
// error message
func (e *PartitionTableCompatibilityError) Error() string {
	return e.err.Error()
}

// Unwrap returns the original error wrapped inside.
func (e *PartitionTableCompatibilityError) Unwrap() error {
	return e.err
}
