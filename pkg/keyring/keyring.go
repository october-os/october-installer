package keyring

import "os/exec"

func UpdateKeyRing() error {
	cmd := exec.Command("pacman", "-Sy", "--noconfirm", "archlinux-keyring")
	if err := cmd.Start(); err != nil {
		return KeyringUpdateError{err: err}
	}
	return nil
}
