package json_parser

import "fmt"

// JsonParsingError represents an error that occured
// trying to parse the json
type JsonParsingError struct {
	Err error
}

// Returns a formatted error message including the underlying
// error message
func (e *JsonParsingError) Error() string {
	return fmt.Sprintf("validation error: error=%v", e.Err)
}

// Returns the error
func (e *JsonParsingError) Unwrap() error {
	return e.Err
}
