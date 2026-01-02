package postinstall

import (
	"fmt"
)

type PostInstallError struct {
	err error
}

func (e PostInstallError) Error() string {
	return fmt.Sprintf("error during post-installation: error=%s", e.err.Error())
}

func (e PostInstallError) Unwrap() error {
	return e.err
}
