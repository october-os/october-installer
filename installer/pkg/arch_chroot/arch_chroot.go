package arch_chroot

import (
	"io"
	"os/exec"
)

const mount_point string = "/mnt"
const shell string = "/bin/bash"

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

func Run(command string) error {
	cmd := exec.Command("arch-chroot", mount_point, shell, "-c", command)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return PipeError{Err: err.Error()}
	}

	err = cmd.Run()
	if err != nil {
		stdErrOutput, _ := io.ReadAll(stderr)
		return ArchChrootError{
			StdErr: string(stdErrOutput),
			Err:    err.Error(),
		}
	}

	return nil
}
