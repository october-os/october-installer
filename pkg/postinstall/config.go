package postinstall

import (
	"fmt"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/user"
)

const octoberConfigRepo string = "https://github.com/october-os/october-config.git"

// Clones the github repo of october-config and installs
// it for each user.
//
// Can return errors of type :
//   - PostInstallError
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

// Goes through all the users and configures
// the config for each user.
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

		if err := arch_chroot.Run(createDirCmd, copyCmd, chownCmd, runSetup); err != nil {
			return err
		}

		copyHyprPreRender := fmt.Sprintf("sudo -u %s cp %s/pre-rendered-templates/colors.conf %s/hypr/base", u.Username, octoberConfigDir, octoberConfigDir)
		copyQuickshellPreRender := fmt.Sprintf("sudo -u %s cp %s/pre-rendered-templates/Theme.qml %s/quickshell/theme", u.Username, octoberConfigDir, octoberConfigDir)
		if err := arch_chroot.Run(copyHyprPreRender, copyQuickshellPreRender); err != nil {
			return err
		}
	}

	return nil
}

// Sets up greetd with tuigreet using sed.
func setupGreetd() error {
	cmd := "sed -i 's/command = \"agreety --cmd \\/bin\\/sh\"/command = \"tuigreet --cmd start-hyprland\"/' /etc/greetd/config.toml"
	if err := arch_chroot.Run(cmd); err != nil {
		return err
	}

	cmd = `echo -e "#%PAM-1.0

auth	required	pam_securetty.so
auth	requisite	pam_nologin.so
auth	include	system-local-login
auth	optional	pam_gnome_keyring.so

account	include	system-local-login
session	include	system-local-login
session	optional	pam_gnome_keyring.so auto_start
" > /etc/pam.d/greetd`

	return arch_chroot.Run(cmd)
}

// Deletes /october-temp from newly installed system.
func removeTemp() error {
	cmd := "rm -rf /october-temp"
	return arch_chroot.Run(cmd)
}

// Creates directory and clones the repo to /october-temp.
func cloneRepoToTemp() error {
	createTmpFolder := "mkdir /october-temp"
	cloneRepo := fmt.Sprintf("git clone %s /october-temp/october-config", octoberConfigRepo)

	return arch_chroot.Run(createTmpFolder, cloneRepo)
}
