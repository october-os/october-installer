package utils

import (
	"os/exec"
)

var CommandExecutorGot []string

type CommandExecutorMock struct {
	*exec.Cmd
}

func (mock CommandExecutorMock) Run() error {
	CommandExecutorGot = append(CommandExecutorGot, mock.Cmd.String())
	return nil
}

func NewCommandExecutorMock(command string, args ...string) ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutorMock{
		Cmd: execCmd,
	}
}
