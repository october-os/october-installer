package fstab

import (
	"os"

	"github.com/october-os/october-installer/pkg/utils"
)

// commandExecutor for mocking exec.Cmd runs.
var commandExecutor = utils.NewCommandExecutor

// fstabPath absolute path to the fstab file.
var fstabPath string = "/mnt/etc/fstab"

// Generates the fstab file with genfstab.
//
// Can return errors of type:
//   - FstabError
func GenerateFstab() error {
	if err := genfstab(); err != nil {
		return FstabError{
			err: err,
		}
	}

	return nil
}

// Generates the fstab file with genfstab.
func genfstab() error {
	command := commandExecutor("genfstab", "-U", "/mnt")
	out, err := command.Output()
	if err != nil {
		return err
	}

	return os.WriteFile(fstabPath, out, 0644)
}
