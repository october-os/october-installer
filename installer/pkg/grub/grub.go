package grub

import (
	"fmt"

	"github.com/arch-couple/arch-couple-installer/pkg/arch_chroot"
)

const esp string = "/boot"
const bootloaderId string = "GRUB"

func InstallGrub() error {
	if err := grubInstall(); err != nil {
		return err
	}

	if err := setUpOsProber(); err != nil {
		return err
	}

	return updateGrubConfig()
}

func updateGrubConfig() error {
	command := "grub-mkconfig -o /boot/grub/grub.cfg"
	return arch_chroot.Run(command)
}

func setUpOsProber() error {
	sedCommand := "sed -i 's/#GRUB_DISABLE_OS_PROBER=false/GRUB_DISABLE_OS_PROBER=false/' /etc/default/grub"
	osProberCommand := "os-prober"
	command := fmt.Sprintf("%s && %s", sedCommand, osProberCommand)
	return arch_chroot.Run(command)
}

func grubInstall() error {
	command := fmt.Sprintf("grub-install --target=x86_64-efi --efi-directory=%s --bootloader-id=%s", esp, bootloaderId)
	return arch_chroot.Run(command)
}
