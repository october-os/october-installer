package hostname

import "fmt"

// HostnameError represents an error that occured
// when trying to set up the network hostname of the
// new system.
type HostnameError struct {
	Err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e HostnameError) Error() string {
	return fmt.Sprintf("Hostname error: error=%s", e.Err.Error())
}

// Unwrap returns the original error wrapped inside
// HostnameError.
func (e HostnameError) Unwrap() error {
	return e.Err
}
