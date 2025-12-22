package core

import (
	"errors"
	"os/exec"
)

// Basic Arch Linux install packages names
const _LINUX_KERNEL string = "linux"
const _BASE_ARCH string = "base"
const _BASE_LINUX_FIRMWARE string = "linux-firmware"

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
			Err: err,
		}
	} else if cpuMicrocode == "" {
		return CoreInstallError{
			Err: errors.New("Unsupported CPU detected. Needs to be an AMD or Intel X86_64 CPU."),
		}
	}

	cmd := exec.Command("pacstrap", "-K", "/mnt", _BASE_ARCH, _LINUX_KERNEL, _BASE_LINUX_FIRMWARE, cpuMicrocode)

	if err := cmd.Run(); err != nil {
		return CoreInstallError{
			Err: err,
		}
	}

	return nil
}
