package postinstall

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// pkg represents a package
// with all its related flags.
type pkg struct {
	name  string
	flags []string
}

// ExtraPackages represents the extra packages to be installed on the system.
// They can come from either:
// - the official Arch Linux repositories
// - the Arch User Repository (AUR)
type ExtraPackages struct {
	OfficialRepositories []string `json:"officialRepositories"`
	AUR                  []string `json:"aur"`
}

// Validates the attributes of an ExtraPackages struct
// Returns a PostInstallError if validation fails
func (e *ExtraPackages) Validate() error {
	whiteSpaceError := PostInstallError{err: errors.New("Extra packages validation: error=white space detected in a package's name")}
	for _, p := range e.OfficialRepositories {
		if strings.Contains(p, " ") {
			return whiteSpaceError
		}
	}

	for _, p := range e.AUR {
		if strings.Contains(p, " ") {
			return whiteSpaceError
		}
	}

	return nil
}

const officialPackagesFilePath string = "/root/postinstall/packages"
const aurPackagesFilePath string = "/root/postinstall/aur"

// Downloads all the packages with a package manager.
// inside the newly installed system.
// It uses pacman or yay if the packages come from the AUR.
func downloadAllPackages(packages []pkg, fromAur bool) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p.name + " ")
	}

	var command string
	if fromAur {
		command = fmt.Sprintf("sudo -u builder yay -S %s --noconfirm", sb.String())
	} else {
		command = fmt.Sprintf("pacman -Sy %s --noconfirm", sb.String())
	}

	return arch_chroot.Run(command)
}

// takes the given file path and read the package
// list inside of it.
//
// List must have this form (flags can be omitted):
//   - [package_name] [flag] [flag]
func getPackageList(path string) ([]pkg, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	var packageList []pkg

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) != 0 && line[0] == '-' {
			splittedLine := strings.Split(line, " ")
			newPackage := pkg{
				name: splittedLine[1],
			}

			if len(splittedLine) > 2 {
				var flags []string
				for i := 2; i < len(splittedLine); i++ {
					flags = append(flags, splittedLine[i])
				}

				newPackage.flags = flags
			}
			packageList = append(packageList, newPackage)
		}
	}

	return packageList, nil
}

// Takes the package slice and parses the flags
// for each packages, then acts accordingly for
// the found flags.
func packageFlagParser(packages []pkg) error {
	var systemdEnablePkgs []string
	var systemdUserEnablePkgs []string

	for _, p := range packages {
		if len(p.flags) == 0 {
			continue
		}

		for _, f := range p.flags {
			if strings.Contains(f, "SYSTEMD_ENABLE") {
				systemdEnablePkgs = append(systemdEnablePkgs, getServiceName(p.name, f))
			} else if strings.Contains(f, "SYSTEMD_USER_ENABLE") {
				systemdUserEnablePkgs = append(systemdUserEnablePkgs, getServiceName(p.name, f))
			}
		}
	}

	if len(systemdEnablePkgs) > 0 {
		if err := systemdEnable(systemdEnablePkgs); err != nil {
			return err
		}
	}

	if len(systemdUserEnablePkgs) > 0 {
		if err := systemdUserEnable(systemdUserEnablePkgs); err != nil {
			return err
		}
	}

	return nil
}

// Checks the flag and returns the real systemd
// service name.
func getServiceName(packageName, flag string) string {
	if strings.Contains(flag, "=") {
		splitString := strings.Split(flag, "=")

		return splitString[1]
	}

	return packageName
}

// Adds the packages to a given file path in this format: "- package\n"
func addPackages(filePath string, packages []string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return err
	}
	defer file.Close()
	for _, p := range packages {
		fmt.Fprintf(file, "- %s\n", p)
	}
	return nil
}
