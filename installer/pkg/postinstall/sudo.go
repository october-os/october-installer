package postinstall

import (
	"fmt"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

const username string = "builder"
const sudoerFilePath string = "/etc/sudoers.d/builder"

func activateBuilderAccount() error {
	createUserCommand := fmt.Sprintf("useradd %s && passwd -d %s", username, username)
	if err := arch_chroot.Run(createUserCommand); err != nil {
		return err
	}

	addingSudoer := fmt.Sprintf("echo '%s ALL=(ALL) NOPASSWD:ALL' > %s", username, sudoerFilePath)
	chmodSudoFile := fmt.Sprintf("chmod 440 %s", sudoerFilePath)

	command := fmt.Sprintf("%s && %s", addingSudoer, chmodSudoFile)
	return arch_chroot.Run(command)
}

func deleteBuilderAccount() error {
	deleteSudoFile := fmt.Sprintf("rm %s", sudoerFilePath)
	deleteUser := fmt.Sprintf("userdel -r %s", username)

	command := fmt.Sprintf("%s && %s", deleteSudoFile, deleteUser)
	return arch_chroot.Run(command)
}
