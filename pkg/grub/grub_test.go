package grub

import (
	"fmt"
	"testing"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestInstallGrub(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := InstallGrub()

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 4, "InstallGrub should run 4 chroot commands")
	assert.Equal(t, grubChrootCommandBuilder(
		"grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=GRUB"),
		mocks.CommandExecutorGot[0], "grub-install should be run first")
	assert.Equal(t, grubChrootCommandBuilder(
		"sed -i 's/#GRUB_DISABLE_OS_PROBER=false/GRUB_DISABLE_OS_PROBER=false/' /etc/default/grub && os-prober"),
		mocks.CommandExecutorGot[1], "os-prober setup should be run second")
	assert.Equal(t, grubChrootCommandBuilder(
		`sed -i 's/GRUB_DISTRIBUTOR="Arch"/GRUB_DISTRIBUTOR="October"/' /etc/default/grub`),
		mocks.CommandExecutorGot[2], "branding should be run third")
	assert.Equal(t, grubChrootCommandBuilder("grub-mkconfig -o /boot/grub/grub.cfg"),
		mocks.CommandExecutorGot[3], "grub-mkconfig should be run last")
}

func TestInstallGrubWithError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	err := InstallGrub()

	assert.Error(t, err)
	assert.ErrorAs(t, err, &GrubError{}, "Error should be wrapped in a GrubError")
	assert.Len(t, mocks.CommandExecutorGot, 1, "only the first command should run before the error")
}

func TestGrubInstall(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := grubInstall()

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(t, grubChrootCommandBuilder(
		"grub-install --target=x86_64-efi --efi-directory=/boot --bootloader-id=GRUB"),
		mocks.CommandExecutorGot[0])
}

func TestSetUpOsProber(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := setUpOsProber()

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(t, grubChrootCommandBuilder(
		"sed -i 's/#GRUB_DISABLE_OS_PROBER=false/GRUB_DISABLE_OS_PROBER=false/' /etc/default/grub && os-prober"),
		mocks.CommandExecutorGot[0])
}

func TestSetUpBranding(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := setUpBranding()

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(t, grubChrootCommandBuilder(
		`sed -i 's/GRUB_DISTRIBUTOR="Arch"/GRUB_DISTRIBUTOR="October"/' /etc/default/grub`),
		mocks.CommandExecutorGot[0])
}

func TestUpdateGrubConfig(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := updateGrubConfig()

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(t, grubChrootCommandBuilder("grub-mkconfig -o /boot/grub/grub.cfg"),
		mocks.CommandExecutorGot[0])
}

func TestSetUpOsProberError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	err := setUpOsProber()

	assert.Error(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
}

func grubChrootCommandBuilder(command string) string {
	return fmt.Sprintf("/usr/bin/arch-chroot /mnt /bin/bash -c %s", command)
}
