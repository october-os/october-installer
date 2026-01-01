package grub

import (
	"fmt"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

const espMountPoint string = "/boot"
const bootloaderId string = "GRUB"

// Installs and sets up Grub on the newly installed system.
//
// Does:
//   - grub-Install
//   - uncommend os-prober line in /etc/default/grub
//   - os-prober
//   - grub-mkconfig
//
// Can return error types:
//   - PipeError
//   - ArchChrootError
func InstallGrub() error {
	if err := grubInstall(); err != nil {
		return err
	}

	if err := setUpOsProber(); err != nil {
		return err
	}

	return updateGrubConfig()
}

// Updates the current Grub config.
//
// Executes:
//
//	grub-mkconfig -o /boot/grub/grub.cfg
func updateGrubConfig() error {
	command := "grub-mkconfig -o /boot/grub/grub.cfg"
	return arch_chroot.Run(command)
}

// Uncomments the os-prober line inside /etc/default/grub
// and runs os-prober.
func setUpOsProber() error {
	sedCommand := "sed -i 's/#GRUB_DISABLE_OS_PROBER=false/GRUB_DISABLE_OS_PROBER=false/' /etc/default/grub"
	osProberCommand := "os-prober"
	command := fmt.Sprintf("%s && %s", sedCommand, osProberCommand)
	return arch_chroot.Run(command)
}

// Runs the Grub installation on the new system.
//
// Executes:
//
//	grub-install...
func grubInstall() error {
	command := fmt.Sprintf(
		"grub-install --target=x86_64-efi --efi-directory=%s --bootloader-id=%s",
		espMountPoint,
		bootloaderId)
	return arch_chroot.Run(command)
}
