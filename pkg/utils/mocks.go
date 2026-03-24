package utils

import (
	"os/exec"
)

// A slice of all the commands that called Run()
var CommandExecutorGot []string

// The struct for mocking CommandExecutor in testing.
type CommandExecutorMock struct {
	*exec.Cmd
}

// Fake Run method that appends a new command to
// CommandExecutorGot.
func (mock CommandExecutorMock) Run() error {
	CommandExecutorGot = append(CommandExecutorGot, mock.Cmd.String())
	return nil
}

// Creates a new mock.
func NewCommandExecutorMock(command string, args ...string) ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutorMock{
		Cmd: execCmd,
	}
}
