package keyring

import "fmt"

// KeyringUpdateError represents an error that occured
// trying to parse the json
type KeyringUpdateError struct {
	err error
}

// Returns a formatted error message including the underlying
// error message
func (e KeyringUpdateError) Error() string {
	return fmt.Sprintf("Error updating keyring: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e KeyringUpdateError) Unwrap() error {
	return e.err
}
