package hostname

import "fmt"

type InvalidHostnameError struct {
	Err error
}

func (e *InvalidHostnameError) Error() string {
	return fmt.Sprintf("invalid hostname error: error=%v", e.Err)
}

func (e *InvalidHostnameError) Unwrap() error {
	return e.Err
}
