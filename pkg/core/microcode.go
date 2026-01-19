package core

import (
	"io"
	"os/exec"
	"strings"
)

// vendor_id field values inside /proc/cpuinfo
const amdId string = "AuthenticAMD"
const intelId string = "GenuineIntel"

// microcode packages name
const amdMicrocode string = "amd-ucode"
const intelMicrocode string = "intel-ucode"

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

	if strings.Contains(string(stdoutBytes), amdId) {
		return amdMicrocode, nil
	} else if strings.Contains(string(stdoutBytes), intelId) {
		return intelMicrocode, nil
	}

	return "", nil
}
