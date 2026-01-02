package postinstall

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

const packageFilePath string = "/root/postinstall/packages"

func downloadAllPackages(packages []string) error {
	var sb strings.Builder
	for _, p := range packages {
		sb.WriteString(p)
		sb.WriteString(" ")
	}

	command := fmt.Sprintf("pacman -S --noconfirm %s", sb.String())
	if err := arch_chroot.Run(command); err != nil {
		return err
	}

	return nil
}

func getPackageList() ([]string, error) {
	fd, err := os.Open(packageFilePath)
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
