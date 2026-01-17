package locale

import "fmt"

// LocaleError represents an error that occured
// when setting up or validating the new install locales.
type LocaleError struct {
	err error
}

// Error returns a formatted error message containing the
// original error message inside.
func (e LocaleError) Error() string {
	return fmt.Sprintf("Error setting up locale: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e LocaleError) Unwrap() error {
	return e.err
}
