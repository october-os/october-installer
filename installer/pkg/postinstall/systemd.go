package postinstall

import (
	"fmt"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

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

// Takes a list of service names that needs to be
// enabled in user systemd and enables them.
//
// By "user systemd" it is meant:
//
//	systemd --user enable [package]
func systemdUserEnable(services []string) error {
	var sb strings.Builder
	for _, s := range services {
		sb.WriteString(s + " ")
	}

	command := fmt.Sprintf("systemctl --user enable %s", sb.String())
	return arch_chroot.Run(command)
}
