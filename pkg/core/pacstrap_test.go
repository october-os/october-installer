package core

import (
	"fmt"
	"testing"

	"github.com/klauspost/cpuid/v2"
	"github.com/october-os/october-installer/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestInstallBasicInstallationValidCpu(t *testing.T) {
	want := fmt.Sprintf("/usr/bin/pacstrap -K /mnt base linux linux-firmware %s", amdMicrocode)
	cpuid.CPU.VendorID = cpuid.AMD
	normalCommandExecutor := commandExecutor
	defer func() {
		commandExecutor = normalCommandExecutor
		utils.T = nil
		utils.Want = ""
	}()

	commandExecutor = utils.NewCommandExecutorMock
	utils.Want = want
	utils.T = t

	err := InstallBasicInstallation()

	assert.Nil(t, err, "Error should be nil with valid installation")
}

func TestInstallBasicInstallationInvalidCpu(t *testing.T) {
	cpuid.CPU.VendorID = cpuid.Apple
	normalCommandExecutor := commandExecutor
	defer func() {
		commandExecutor = normalCommandExecutor
		utils.T = nil
	}()

	commandExecutor = utils.NewCommandExecutorMock
	utils.T = t

	err := InstallBasicInstallation()

	assert.NotNil(t, err, "Error should be not nil with valid installation")
	assert.ErrorAs(t, err, &CoreInstallError{}, "Error should be CoreInstallerError on error")
}
