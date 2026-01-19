package core

import "fmt"

// CoreInstallError represents an error that occured
// during the installation of the core Arch and Linux system.
type CoreInstallError struct {
	err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e CoreInstallError) Error() string {
	return fmt.Sprintf("Error during core installation: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e CoreInstallError) Unwrap() error {
	return e.err
}
