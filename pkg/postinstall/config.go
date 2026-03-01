package postinstall

import (
	"fmt"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/user"
)

const octoberConfigRepo string = "https://github.com/october-os/october-config.git"

func AddConfigForUsers(users []user.User) error {
	if err := cloneRepoToTemp(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := setupConfigForUsers(users); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}

func setupConfigForUsers(users []user.User) error {
	for _, u := range users {
		createDirCmd := fmt.Sprintf("mkdir %s/.config", u.Homepath)
		copyCmd := fmt.Sprintf("cp /tmp/october-config %s/.config", u.Homepath)
		chownCmd := fmt.Sprintf("chown -R %s:%s %s/.config", u.Username, u.Username, u.Homepath)
		runSetup := fmt.Sprintf("sudo -u %s %s/.config/october-config/scripts/setup.sh", u.Username, u.Homepath)

		cmd := fmt.Sprintf("%s && %s && %s && %s", createDirCmd, copyCmd, chownCmd, runSetup)

		if err := arch_chroot.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}

func cloneRepoToTemp() error {
	cmd := fmt.Sprintf("git clone %s /tmp/october-config", octoberConfigRepo)

	return arch_chroot.Run(cmd)
}
