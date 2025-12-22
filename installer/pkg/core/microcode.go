package core

import (
	"io"
	"os/exec"
	"strings"
)

const AMD_ID string = "AuthenticAMD"
const INTEL_ID string = "GenuineIntel"

const AMD_MICROCODE string = "amd-ucode"
const INTEL_MICROCODE string = "intel-ucode"

func GettCpuMicroCode() (string, error) {
	cmd := exec.Command("/bin/bash", "-c", "cat /proc/cpuinfo | grep 'vendor_id'")
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	stdoutBytes, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return "", err
	}

	if err := cmd.Wait(); err != nil {
		return "", err
	}

	if strings.Contains(string(stdoutBytes), AMD_ID) {
		return AMD_MICROCODE, nil
	} else if strings.Contains(string(stdoutBytes), INTEL_ID) {
		return INTEL_MICROCODE, nil
	}

	return "", nil
}
