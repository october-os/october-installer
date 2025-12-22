package core

import "fmt"

type CoreInstallError struct {
	Err error
}

func (e CoreInstallError) Error() string {
	return fmt.Sprintf("Core installer error: error=%s", e.Err.Error())
}

func (e CoreInstallError) Unwrap() error {
	return e.Err
}
