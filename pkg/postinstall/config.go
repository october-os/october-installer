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

	if err := removeTemp(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := setupGreetd(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}

func setupConfigForUsers(users []user.User) error {
	for _, u := range users {
		if u.Homepath == "" {
			u.Homepath = fmt.Sprintf("/home/%s", u.Username)
		}

		octoberConfigDir := fmt.Sprintf("%s/.config/october-config", u.Homepath)

		// Creating the whole path since we need those folders
		createDirCmd := fmt.Sprintf("mkdir -p %s/.config/wal/templates", u.Homepath)
		copyCmd := fmt.Sprintf("cp -r /october-temp/october-config %s/.config", u.Homepath)
		chownCmd := fmt.Sprintf("chown -R %s:%s %s/.config", u.Username, u.Username, u.Homepath)
		runSetup := fmt.Sprintf("sudo -u %s %s/scripts/setup.sh", u.Username, octoberConfigDir)

		cmd := fmt.Sprintf("%s && %s && %s && %s", createDirCmd, copyCmd, chownCmd, runSetup)

		if err := arch_chroot.Run(cmd); err != nil {
			return err
		}

		copyHyprPreRender := fmt.Sprintf("sudo -u %s cp %s/pre-rendered-templates/colors.conf %s/hypr/base", u.Username, octoberConfigDir, octoberConfigDir)
		copyQuickshellPreRender := fmt.Sprintf("sudo -u %s cp %s/pre-rendered-templates/Theme.qml %s/quickshell/theme", u.Username, octoberConfigDir, octoberConfigDir)
		cmd = fmt.Sprintf("%s && %s", copyHyprPreRender, copyQuickshellPreRender)
		if err := arch_chroot.Run(cmd); err != nil {
			return err
		}
	}

	return nil
}

func setupGreetd() error {
	cmd := "sed -i 's/command = \"agreety --cmd \\/bin\\/sh\"/command = \"tuigreet --cmd start-hyprland\"/' /etc/greetd/config.toml"
	return arch_chroot.Run(cmd)
}

func removeTemp() error {
	cmd := "rm -rf /october-temp"
	return arch_chroot.Run(cmd)
}

func cloneRepoToTemp() error {
	createTmpFolder := "mkdir /october-temp"
	cloneRepo := fmt.Sprintf("git clone %s /october-temp/october-config", octoberConfigRepo)

	cmd := fmt.Sprintf("%s && %s", createTmpFolder, cloneRepo)

	return arch_chroot.Run(cmd)
}
