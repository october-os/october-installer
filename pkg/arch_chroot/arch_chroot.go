package arch_chroot

import (
	"io"
	"os/exec"
	"strings"
)

// mountPoint is the mount point of the system to chroot into
const mountPoint string = "/mnt"

// shell is the shell that will be used to execute chroot commands
const shell string = "/bin/bash"

// Executes the command in a shell using arch-chroot.
//
// It executes: arch-chroot [mount_point] [shell] -c [command]
//
// It can return types of errors:
//   - ArchChrootError
func Run(commands... string) error {
	var command strings.Builder
	if len(commands) == 1 {
		command.WriteString(commands[0])
	} else {
		for _, cmd := range commands {
			command.WriteString(cmd)
			command.WriteString(" && ")
		}
	}

	cmd := exec.Command("arch-chroot", mountPoint, shell, "-c", command.String())
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return ArchChrootError{err: err}
	}

	err = cmd.Run()
	if err != nil {
		stdErrOutput, _ := io.ReadAll(stderr)
		return ArchChrootError{
			command: command.String(),
			stdErr:  string(stdErrOutput),
			err:     err,
		}
	}

	return nil
}
