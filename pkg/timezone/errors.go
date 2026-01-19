package timezone

import "fmt"

// TimezoneError represents an error that occured
// when setting up or validating the new install timezones.
type TimezoneError struct {
	err error
}

// Returns a formatted error message including the original
// error message.
func (e TimezoneError) Error() string {
	return fmt.Sprintf("Error setting timezone: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e TimezoneError) Unwrap() error {
	return e.err
}
