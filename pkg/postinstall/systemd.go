package postinstall

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

const systemdServicesFilePath string = "/root/postinstall/services"

func getSystemdServices() ([]string, error) {
	f, err := os.Open(systemdServicesFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var services []string

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			continue
		}

		if line[0] == '-' {
			serviceName := strings.TrimPrefix(line, "- ")
			services = append(services, serviceName)
		}
	}

	return services, nil
}

// Takes a list of service names that needs to be
// enabled in systemd and enables them.
func systemdEnable(services []string) error {
	var sb strings.Builder
	for _, s := range services {
		sb.WriteString(s + " ")
	}

	command := fmt.Sprintf("systemctl enable %s", sb.String())
	return arch_chroot.Run(command)
}
