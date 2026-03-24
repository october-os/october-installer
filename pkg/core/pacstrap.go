package core

import (
	"github.com/october-os/october-installer/pkg/utils"
)

var commandExecutor = utils.NewCommandExecutor

// Basic Arch Linux install packages names
const linuxKernel string = "linux"
const baseArch string = "base"
const baseLinuxFirmware string = "linux-firmware"

// Installs a basic Arch Linux installation on the drive
// mounted on /mnt using pacstrap. Detects and installs the CPU
// microcode for the current CPU too.
//
// Can return errors of type:
//   - CoreInstallError
func InstallBasicInstallation() error {
	cpuMicrocode, err := getCpuMicroCode()
	if err != nil {
		return CoreInstallError{
			err: err,
		}
	}

	cmd := commandExecutor(
		"pacstrap",
		"-K",
		"/mnt",
		baseArch,
		linuxKernel,
		baseLinuxFirmware,
		cpuMicrocode)

	if err := cmd.Run(); err != nil {
		return CoreInstallError{
			err: err,
		}
	}

	return nil
}
