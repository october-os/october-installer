package postinstall

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

const packageFilePath string = "/root/postinstall/packages"
const aurFilePath string = "/root/postinstall/aur"

func downloadAllPackages(packages []string, aur bool) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p)
		sb.WriteString(" ")
	}

	var command string
	if aur {
		command = fmt.Sprintf("sudo -u builder yay -S %s --noconfirm", sb.String())
	} else {
		command = fmt.Sprintf("pacman -S --noconfirm %s", sb.String())
	}

	return arch_chroot.Run(command)
}

func getPackageList(path string) ([]string, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	scanner := bufio.NewScanner(fd)
	var packageList []string

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) != 0 && line[0] == '-' {
			packageList = append(packageList, strings.Split(line, " ")[1])
		}
	}

	return packageList, nil
}
