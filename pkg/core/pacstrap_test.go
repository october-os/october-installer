package core

import (
	"fmt"
	"testing"

	"github.com/klauspost/cpuid/v2"
	"github.com/stretchr/testify/assert"
)

func TestInstallBasicInstallationValidCpu(t *testing.T) {
	want := fmt.Sprintf("/usr/bin/pacstrap -K /mnt base linux linux-firmware %s", amdMicrocode)
	cpuid.CPU.VendorID = cpuid.AMD

	currentlyTesting = true
	defer func() {
		currentlyTesting = false
		testCmd = nil
	}()

	err := InstallBasicInstallation()

	assert.Nil(t, err, "Error should be nil with valid installation")
	assert.Equal(t, want, testCmd.String())
}

func TestInstallBasicInstallationInvalidCpu(t *testing.T) {
	cpuid.CPU.VendorID = cpuid.Apple

	currentlyTesting = true
	defer func() {
		currentlyTesting = false
		testCmd = nil
	}()

	err := InstallBasicInstallation()

	assert.NotNil(t, err, "Error should be not nil with valid installation")
	assert.ErrorAs(t, err, &CoreInstallError{}, "Error should be CoreInstallerError on error")
	assert.Nil(t, testCmd, "Test cmd should be nil with invalid installation")
}
