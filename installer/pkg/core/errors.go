package core

import "fmt"

// CoreInstallError represents an error that occured
// during the instrallation of the core Arch and Linux system.
type CoreInstallError struct {
	Err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e CoreInstallError) Error() string {
	return fmt.Sprintf("Core installer error: error=%s", e.Err.Error())
}

// Unwrap returns the original error wrapped inside
// CoreInstallError.
func (e CoreInstallError) Unwrap() error {
	return e.Err
}
