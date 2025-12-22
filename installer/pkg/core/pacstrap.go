package core

import (
	"errors"
	"os/exec"
)

const LINUX_KERNEL string = "linux"
const BASE_ARCH string = "base"
const BASE_LINUX_FIRMWARE string = "linux-firmware"

func InstallBasicInstallation() error {
	cpuMicrocode, err := GettCpuMicroCode()
	if err != nil {
		return CoreInstallError{
			Err: err,
		}
	} else if cpuMicrocode == "" {
		return CoreInstallError{
			Err: errors.New("Unsupported CPU detected. Needs to be an AMD or Intel X86_64 CPU."),
		}
	}

	cmd := exec.Command("pacstrap", "-K", "/mnt", BASE_ARCH, LINUX_KERNEL, BASE_LINUX_FIRMWARE, cpuMicrocode)

	if err := cmd.Run(); err != nil {
		return CoreInstallError{
			Err: err,
		}
	}

	return nil
}
