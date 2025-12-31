package mirrors

import "fmt"

type MirrorListError struct {
	err error
}

func (e MirrorListError) Error() string {
	return fmt.Sprintf("mirrorlist error: error=%s", e.err.Error())
}

func (e MirrorListError) Unwrap() error {
	return e.err
}
