package postinstall

import (
	"fmt"
)

// PostInstallError represents an error that occured
// during the post-installation steps of the
// installation.
type PostInstallError struct {
	err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e PostInstallError) Error() string {
	return fmt.Sprintf("error during post-installation: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e PostInstallError) Unwrap() error {
	return e.err
}
