package arch_chroot

type ArchChrootError struct {
	StdErr string
	Err    string
}

type PipeError struct {
	Err string
}

func (e ArchChrootError) Error() string {
	return e.Err
}

func (e PipeError) Error() string {
	return e.Err
}
