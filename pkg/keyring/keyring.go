package keyring

import (
	"github.com/october-os/october-installer/pkg/utils"
)

// commandExecutor is the command executor for the keyring
// package.
var commandExecutor = utils.NewCommandExecutor

// UpdateKeyRing updates the current Arch Linux keyring of the ISO
// so it can install October even if it is out of date.
//
// Can return error of type:
// 	- KeyringUpdateError
func UpdateKeyRing() error {
	cmd := commandExecutor("pacman", "-Sy", "--noconfirm", "archlinux-keyring")
	if err := cmd.Run(); err != nil {
		return KeyringUpdateError{err: err}
	}
	return nil
}
