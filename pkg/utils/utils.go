package utils

import "os/exec"

type ICommandExecutor interface {
	Run() error
}

type CommandExecutor struct {
	*exec.Cmd
}

func (c CommandExecutor) Run() error {
	return c.Cmd.Run()
}

func NewCommandExecutor(command string, args ...string) ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutor{execCmd}
}
