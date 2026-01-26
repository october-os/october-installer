package postinstall

import (
	"io"
	"os"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Gets the list of packages that need to be installed
// with pacman, installs them, then configure them
// if needed.
//
// Can return errors of types :
//   - PostInstallError
func InstallOfficialPackages() error {
	packages, err := getPackageList(officialPackagesFilePath)
	if err != nil {
		return PostInstallError{err: err}
	}

	if err := downloadAllPackages(packages, false); err != nil {
		return PostInstallError{err: err}
	}

	if err := packageFlagParser(packages); err != nil {
		return PostInstallError{err: err}
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
		return PostInstallError{err: err}
	}

	if err := downloadAllPackages(packages, true); err != nil {
		return PostInstallError{err: err}
	}

	if err := deleteBuilderAccount(); err != nil {
		return PostInstallError{err: err}
	}

	if err := packageFlagParser(packages); err != nil {
		return PostInstallError{err: err}
	}

	return nil
}

// Enables the multilib package repository in pacman.
//
// Can return errors of types :
//   - PostInstallError
func EnableMultilibRepo() error {
	command := "sed -i -e '/#\\[multilib\\]/,+1s/^#//' /etc/pacman.conf"
	if err := arch_chroot.Run(command); err != nil {
		return PostInstallError{err: err}
	}

	return nil
}

// Enables the wheel group in sudo.
//
// Can return errors of types:
//   - PostInstallError
func SetupSudo() error {
	if err := addWheelGroup(); err != nil {
		return PostInstallError{err: err}
	}

	return nil
}

// Copies the /etc/os-release and /etc/lsb-release files from the ISO to the installed system
// for branding purposes.
//
// Can return errors of types:
//   - PostInstallError
func SetupBranding() error {
	isoFile, err := os.Open("/etc/os-release")
	if err != nil {
		return PostInstallError{err: err}
	}
	systemFile, err := os.Create("/mnt/etc/os-release")
	if err != nil {
		return PostInstallError{err: err}
	}
	if _, err := io.Copy(systemFile, isoFile); err != nil {
		return PostInstallError{err: err}
	}

	isoFile, err = os.Open("/etc/lsb-release")
	if err != nil {
		return PostInstallError{err: err}
	}
	systemFile, err = os.Create("/mnt/etc/lsb-release")
	if err != nil {
		return PostInstallError{err: err}
	}
	if _, err := io.Copy(systemFile, isoFile); err != nil {
		return PostInstallError{err: err}
	}

	return nil
}
