package mirrors

import "fmt"

// MirrorListError represents an error that occured
// during setting or validating mirrors
type MirrorListError struct {
	err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e MirrorListError) Error() string {
	return fmt.Sprintf("Error setting mirrorlist: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e MirrorListError) Unwrap() error {
	return e.err
}
