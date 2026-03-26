package utils

import "os/exec"

// Interface for exec.Cmd.Run(). Makes
// it mockable and testable.
type ICommandExecutor interface {
	Run() error
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

// Creates a new CommandExecutor.
func NewCommandExecutor(command string, args ...string) ICommandExecutor {
	execCmd := exec.Command(command, args...)
	return CommandExecutor{execCmd}
}
