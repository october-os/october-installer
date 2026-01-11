package postinstall

import "github.com/october-os/october-installer/pkg/arch_chroot"

// Gets the list of packages that need to be installed
// with pacman, installs them, then configure them
// if needed.
//
// Can return errors of types :
//   - PostInstallError
func InstallOfficialPackages() error {
	packages, err := getPackageList(officialPackagesFilePath)
	if err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := downloadAllPackages(packages, false); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := packageFlagParser(packages); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}

// Gets the list of packages that need to be installed
// with yay, installs yay and them, then configure them
// if needed.
//
// Can return errors of types :
//   - PostInstallError
func InstallAurHelperAndPackages() error {
	if err := activateBuilderAccount(); err != nil {
		return err
	}

	if err := installYay(); err != nil {
		return err
	}

	packages, err := getPackageList(aurPackagesFilePath)
	if err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := downloadAllPackages(packages, true); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := deleteBuilderAccount(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	if err := packageFlagParser(packages); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}

// Enables the multilib package repository in pacman.
//
// Can return errors of types :
//   - PostInstallError
func EnableMultilibRepo() error {
	command := "sed -i -e '/#\\[multilib\\]/,+1s/^#//' /etc/pacman.conf"
	return arch_chroot.Run(command)
}

// Enables the wheel group in sudo.
//
// Can return errors of types:
//   - PostInstallError
func SetupSudo() error {
	if err := addWheelGroup(); err != nil {
		return PostInstallError{
			err: err,
		}
	}

	return nil
}
