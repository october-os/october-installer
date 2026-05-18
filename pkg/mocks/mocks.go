package mocks

import (
	"errors"
	"io"
	"os/exec"

	"github.com/october-os/october-installer/pkg/utils"
)

// CommandExecutorGot slice of all the commands that called Run()
var CommandExecutorGot []string

// OutputReturn string to be returned by Output()
// when called
var OutputReturn string

// ReturnError boolean to indicate if methods need
// to fail and return an error
var ReturnError bool = false

// CommandExecutorMock struct for mocking CommandExecutor in testing.
// Implements ICommandExecutor.
type CommandExecutorMock struct {
	*exec.Cmd
}

// Run is a mock method that appends a new command to
// CommandExecutorGot.
func (mock CommandExecutorMock) Run() error {
	CommandExecutorGot = append(CommandExecutorGot, mock.String())

	if ReturnError {
		return errors.New("mock error")
	}
	return nil
}

// Output fakes Output() method run and appends
// the command in CommandExecutorGot.
func (mock CommandExecutorMock) Output() ([]byte, error) {
	CommandExecutorGot = append(CommandExecutorGot, mock.String())

	if ReturnError {
		return nil, errors.New("mock error")
	}
	return []byte(OutputReturn), nil
}

func (mock CommandExecutorMock) GetStdErrPipe() (io.ReadCloser, error) {
	return mock.Cmd.StderrPipe()
}

// NewCommandExecutorMock creates a new mock.
func NewCommandExecutorMock(command string, args ...string) utils.ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutorMock{
		Cmd: execCmd,
	}
}
