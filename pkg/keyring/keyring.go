package keyring

import (
	"github.com/october-os/october-installer/pkg/utils"
)

var commandExecutor = utils.NewCommandExecutor

func UpdateKeyRing() error {
	cmd := commandExecutor("pacman", "-Sy", "--noconfirm", "archlinux-keyring")
	if err := cmd.Run(); err != nil {
		return KeyringUpdateError{err: err}
	}
	return nil
}
