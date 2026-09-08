package locale

import (
	"fmt"
	"testing"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

const testdataPath string = "testdata/locale.gen"

func TestValidateLocaleValid(t *testing.T) {
	originalFilePath := filepath
	filepath = testdataPath
	defer func() {
		filepath = originalFilePath
	}()

	err := ValidateLocale("fr_FR.UTF-8")

	assert.NoError(t, err)
}

func TestValidateLocaleInvalid(t *testing.T) {
	originalFilePath := filepath
	filepath = testdataPath
	defer func() {
		filepath = originalFilePath
	}()

	err := ValidateLocale("xx_XX.UTF-8")

	assert.Error(t, err)
	assert.ErrorAs(t, err, &LocaleError{})
}

func TestValidateLocaleNonUTF8(t *testing.T) {
	originalFilePath := filepath
	filepath = testdataPath
	defer func() {
		filepath = originalFilePath
	}()

	err := ValidateLocale("en_US")

	assert.Error(t, err)
	assert.ErrorAs(t, err, &LocaleError{})
}

func TestGenerateLocales(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := GenerateLocales("en_US.UTF-8")

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder(
			"sed -i 's/#en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && "+
				"echo LANG=en_US.UTF-8 > /etc/locale.conf && locale-gen",
		),
		mocks.CommandExecutorGot[0])
	assert.Contains(
		t,
		mocks.CommandExecutorGot[0],
		"s/#en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8")
}

func TestGenerateLocalesError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	err := GenerateLocales("en_US.UTF-8")

	assert.Error(t, err)
	assert.ErrorAs(t, err, &LocaleError{})
}

func archChrootCommandBuilder(command string) string {
	return fmt.Sprintf("/usr/bin/arch-chroot /mnt /bin/bash -c %s", command)
}
