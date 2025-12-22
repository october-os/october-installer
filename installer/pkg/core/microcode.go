package core

import (
	"io"
	"os/exec"
	"strings"
)

// vendor_id field values inside /proc/cpuinfo
const _AMD_ID string = "AuthenticAMD"
const _INTEL_ID string = "GenuineIntel"

// microcode packages name
const _AMD_MICROCODE string = "amd-ucode"
const _INTEL_MICROCODE string = "intel-ucode"

// Checks the vendor_id of all CPUs and returns the
// corresponding microcode package that has to be installed.
//
// It gets the vendor id by executing:
//
//	cat /proc/cpuinfo | grep 'vendor_id'
func getCpuMicroCode() (string, error) {
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

	if strings.Contains(string(stdoutBytes), _AMD_ID) {
		return _AMD_MICROCODE, nil
	} else if strings.Contains(string(stdoutBytes), _INTEL_ID) {
		return _INTEL_MICROCODE, nil
	}

	return "", nil
}
