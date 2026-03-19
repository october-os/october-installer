package core

import (
	"errors"

	"github.com/klauspost/cpuid/v2"
)

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
	switch cpuid.CPU.VendorID {
	case cpuid.Intel:
		return intelMicrocode, nil
	case cpuid.AMD:
		return amdMicrocode, nil
	default:
		return "", CoreInstallError{errors.New("CPU is not supported")}
	}
}
