package fstab

import "fmt"

// FstabError represents an error that occured
// during generating the file or other
// actions made on it.
type FstabError struct {
	err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e FstabError) Error() string {
	return fmt.Sprintf("Error while configuring fstab: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e FstabError) Unwrap() error {
	return e.err
}
