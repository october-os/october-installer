package fstab

import (
	"os"
	"os/exec"
)

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
	out, err := exec.Command("genfstab", "-U", "/mnt").Output()
	if err != nil {
		return err
	}

	return os.WriteFile("/mnt/etc/fstab", out, 0644)
}
