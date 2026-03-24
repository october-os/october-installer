package utils

import (
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

var Want string
var T *testing.T

type CommandExecutorMock struct {
	*exec.Cmd
}

func (mock CommandExecutorMock) Run() error {
	assert.Equal(T, Want, mock.Cmd.String())
	return nil
}

func NewCommandExecutorMock(command string, args ...string) ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutorMock{
		Cmd: execCmd,
	}
}
