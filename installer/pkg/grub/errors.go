package grub

import "fmt"

// GrubError represents an error that occured
// during the installation and setup of GRUB bootloader
type GrubError struct {
	err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e *GrubError) Error() string {
	return fmt.Sprintf("Error setting up GRUB: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e *GrubError) Unwrap() error {
	return e.err
}
