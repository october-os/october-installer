package user

import "fmt"

type PasswordError struct {
	message string
	err     error
}

func (e PasswordError) Error() string {
	return fmt.Sprintf("password error: %q, error=%v", e.message, e.err.Error())
}

func (e PasswordError) Unwrap() error {
	return e.err
}
