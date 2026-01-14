package core

import (
	"errors"
	"os/exec"
)

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
		return &CoreInstallError{
			err: err,
		}
	} else if cpuMicrocode == "" {
		return &CoreInstallError{
			err: errors.New("Unsupported CPU detected. Needs to be an AMD or Intel X86_64 CPU."),
		}
	}

	cmd := exec.Command("pacstrap", "-K", "/mnt", baseArch, linuxKernel, baseLinuxFirmware, cpuMicrocode)

	if err := cmd.Run(); err != nil {
		return &CoreInstallError{
			err: err,
		}
	}

	return nil
}
