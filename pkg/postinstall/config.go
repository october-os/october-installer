package postinstall

import (
	"fmt"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/user"
)

const octoberConfigRepo string = "https://github.com/october-os/october-config.git"

func copyRepoToUser(users []user.User) error {
	for _, u := range users {
		cmd := fmt.Sprintf("cp /tmp/october-config %s", u.Homepath+"/.config")

		if err := arch_chroot.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}

func cloneRepoToTemp() error {
	cmd := fmt.Sprintf("git clone %s /tmp", octoberConfigRepo)

	return arch_chroot.Run(cmd)
}
