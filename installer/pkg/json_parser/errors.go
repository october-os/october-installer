package json_parser

import "fmt"

// JsonParsingError represents an error that occured
// trying to parse the json
type JsonParsingError struct {
	err error
}

// Returns a formatted error message including the underlying
// error message
func (e JsonParsingError) Error() string {
	return fmt.Sprintf("Error validating JSON: error=%s", e.err.Error())
}

// Unwrap returns the original error wrapped inside.
func (e JsonParsingError) Unwrap() error {
	return e.err
}
