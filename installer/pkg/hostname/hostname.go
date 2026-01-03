package hostname

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Sets the network hostname for the newly
// installed system. It sets it inside /etc/hostname.
//
// Can return errors of types:
//   - PipeError
//   - ArchChrootError
func SetHostname(hostname string) error {
	command := fmt.Sprintf("echo %s > /etc/hostname", hostname)
	return arch_chroot.Run(command)
}

// Checks if the given hostname is RFC1178 complient.
//
// For more information: https://wiki.archlinux.org/title/Installation_guide#Network_configuration
//
// Can return errors of types:
//   - HostnameError
func ValidateHostname(hostname string) error {
	var valid bool = true

	if len(hostname) < 1 || len(hostname) > 63 {
		valid = false
	} else if hostname[0] == '-' {
		valid = false
	} else if !charCheck(hostname) {
		valid = false
	}

	if !valid {
		return HostnameError{
			Err: errors.New("Invalid hostname. Must be RFC1178 complient"),
		}
	}

	return nil
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
