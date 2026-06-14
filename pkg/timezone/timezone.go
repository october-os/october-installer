package timezone

import (
	"errors"
	"fmt"
	"os/exec"
	"slices"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Sets timezone up inside the new install.
//
// Can return error types:
//   - TimezoneError
func SetTime(timezone string) error {
	command := fmt.Sprintf("ln -sf /usr/share/zoneinfo/%s /etc/localtime", timezone)
	if err := arch_chroot.Run(command); err != nil {
		return TimezoneError{err: err}
	}
	return nil
}

// Sets up hardware clock to generate /etc/adjtime.
//
// Runs the following command in arch-chroot:
//
//	hwclock --systohc
//
// Can return error types:
//   - TimezoneError
func SetHwClock() error {
	command := "hwclock --systohc"
	if err := arch_chroot.Run(command); err != nil {
		return TimezoneError{err: err}
	}
	return nil
}

// Checks if the given timezone is a valid.
//
// Can return error types:
//   - TimezoneError
func ValidateTimezone(timezone string) error {
	timezones, err := getAllTimezones()
	if err != nil {
		return TimezoneError{err: err}
	}

	if _, found := slices.BinarySearch(timezones, timezone); !found {
		return TimezoneError{
			err: errors.New("Invalid timezone"),
		}
	}

	return nil
}

// Gets all the timezones from STDOUT and returns them
// in an array of string.
//
// It executes:
//
//	timedatectl list-timezones
func getAllTimezones() ([]string, error) {
	cmd := exec.Command("timedatectl", "list-timezones")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return strings.Split(string(stdout), "\n"), nil
}
