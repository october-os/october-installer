package user

import "fmt"

type UserInstantiationError struct {
	err error
}

func (e UserInstantiationError) Error() string {
	return fmt.Sprintf("failed to instantiate user: error=%v", e.err.Error())
}

func (e UserInstantiationError) Unwrap() error {
	return e.err
}
