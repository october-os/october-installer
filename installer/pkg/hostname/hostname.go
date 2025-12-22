package hostname

import (
	"errors"
	"fmt"
	"unicode"

	"github.com/arch-couple/arch-couple-installer/pkg/arch_chroot"
)

func SetHostname(hostname string) error {
	if !isRFC1178Complient(hostname) {
		return HostnameError{
			Err: errors.New("Invalid hostname"),
		}
	}

	command := fmt.Sprintf("echo %s > /etc/hostname", hostname)

	return arch_chroot.Run(command)
}

func isRFC1178Complient(hostname string) bool {
	if len(hostname) < 1 || len(hostname) > 63 {
		return false
	} else if hostname[0] == '-' {
		return false
	} else if !isLowercaseAndContainsNoWhitespaces(hostname) {
		return false
	}

	return true
}

func isLowercaseAndContainsNoWhitespaces(s string) bool {
	for _, r := range s {
		if !unicode.IsLower(r) && unicode.IsLetter(r) || r == ' ' {
			return false
		}
	}

	return true
}
