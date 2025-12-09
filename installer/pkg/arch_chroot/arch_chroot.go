package arch_chroot

import (
	"io"
	"os/exec"
)

// mountPoint is the mount point of the system to chroot into
const mountPoint string = "/mnt"

// shell is the shell that will be used to execute chroot commands
const shell string = "/bin/bash"

// Executes the command in a shell using arch-chroot.
//
// It executes: arch-chroot [mount_point] [shell] -c [command]
//
// It can return two types of errors:
//   - PipeError: When it failed to pipe STDERR
//   - ArchChrootError: When the command ran with arch-chroot failed.
func Run(command string) error {
	cmd := exec.Command("arch-chroot", mountPoint, shell, "-c", command)
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return PipeError{Err: err}
	}

	err = cmd.Run()
	if err != nil {
		stdErrOutput, _ := io.ReadAll(stderr)
		return ArchChrootError{
			StdErr: string(stdErrOutput),
			Err:    err,
		}
	}

	return nil
}
