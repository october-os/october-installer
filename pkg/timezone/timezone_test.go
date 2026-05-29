package timezone

import (
	"fmt"
	"testing"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSetTime(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	timezone := "America/Montreal"
	err := SetTime(timezone)

	assert.Nil(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder("ln -sf /usr/share/zoneinfo/America/Montreal /etc/localtime"),
		mocks.CommandExecutorGot[0])
}

func TestSetTimeError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	timezone := "invalid/tz"
	err := SetTime(timezone)

	assert.NotNil(t, err)
}

func archChrootCommandBuilder(command string) string {
	return fmt.Sprintf("/usr/bin/arch-chroot /mnt /bin/bash -c %s", command)
}
