package fstab

import "os/exec"

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
	cmd := exec.Command("/bin/bash", "-c", "genfstab -U /mnt >> /mnt/etc/fstab")
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
