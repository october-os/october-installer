package utils

import (
	"io"
	"os/exec"
)

// Interface for exec.Cmd.Run(). Makes
// it mockable and testable.
type ICommandExecutor interface {
	Run() error
	Output() ([]byte, error)
	GetStdErrPipe() (io.ReadCloser, error)
}

// CommandExecutor struct that implements
// ICommandExecutor and **runs commands
// bare metal**
type CommandExecutor struct {
	*exec.Cmd
}

// Run method that calls exec.Cmd.Run().
func (c CommandExecutor) Run() error {
	return c.Cmd.Run()
}

// Output runs the Output() method
// of an exec.Cmd.
func (c CommandExecutor) Output() ([]byte, error) {
	return c.Cmd.Output()
}

func (c CommandExecutor) GetStdErrPipe() (io.ReadCloser, error) {
	return c.Cmd.StderrPipe()
}

// Creates a new CommandExecutor.
func NewCommandExecutor(command string, args ...string) ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutor{execCmd}
}
