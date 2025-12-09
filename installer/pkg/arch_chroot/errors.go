package arch_chroot

import "fmt"

type ArchChrootError struct {
	StdErr string
	Err    error
}

func (e ArchChrootError) Error() string {
	return fmt.Sprintf("STDERR: %s, Error: %s", e.StdErr, e.Err.Error())
}

func (e ArchChrootError) Unwrap() error {
	return e.Err
}

type PipeError struct {
	Err error
}

func (e PipeError) Error() string {
	return e.Err.Error()
}

func (e PipeError) Unwrap() error {
	return e.Err
}
