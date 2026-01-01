package hostname

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Checks the validity and Sets the network hostname
// for the newly installed system. It sets it inside /etc/hostname.
//
// Can return errors of types:
//   - HostnameError
//   - PipeError
//   - ArchChrootError
func SetHostname(hostname string) error {
	if !isRFC1178Complient(hostname) {
		return HostnameError{
			Err: errors.New("Invalid hostname"),
		}
	}

	command := fmt.Sprintf("echo %s > /etc/hostname", hostname)

	return arch_chroot.Run(command)
}

// Checks if the given hostname is RFC1178 complient.
//
// For more information: https://wiki.archlinux.org/title/Installation_guide#Network_configuration
func isRFC1178Complient(hostname string) bool {
	if len(hostname) < 1 || len(hostname) > 63 {
		return false
	} else if hostname[0] == '-' {
		return false
	} else if !charCheck(hostname) {
		return false
	}

	return true
}

// Checks if the given string is all lowercase and doesn't
// contain any whitespaces.
func charCheck(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) || r == ' ' {
			return false
		}
	}

	return true
}
