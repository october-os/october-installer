package utils

import (
	"errors"
	"os/exec"
)

// CommandExecutorGot slice of all the commands that called Run()
var CommandExecutorGot []string

var OutputMockReturn string

var ReturnMockError bool = false

// CommandExecutorMock struct for mocking CommandExecutor in testing.
// Implements ICommandExecutor.
type CommandExecutorMock struct {
	*exec.Cmd
}

// Run mock method that appends a new command to
// CommandExecutorGot.
func (mock CommandExecutorMock) Run() error {
	CommandExecutorGot = append(CommandExecutorGot, mock.String())

	if ReturnMockError {
		return errors.New("Mock error")
	}
	return nil
}

// Output fakes Output() method run and appends
// the command in CommandExecutorGot.
func (mock CommandExecutorMock) Output() ([]byte, error) {
	CommandExecutorGot = append(CommandExecutorGot, mock.String())

	if ReturnMockError {
		return nil, errors.New("Mock error")
	}
	return []byte(OutputMockReturn), nil
}

// NewCommandExecutorMock creates a new mock.
func NewCommandExecutorMock(command string, args ...string) ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutorMock{
		Cmd: execCmd,
	}
}
