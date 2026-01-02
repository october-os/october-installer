package timezone

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"slices"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Sets timezone up inside the new install.
//
// Can return error types:
//   - PipeError
//   - ArchChrootError
func SetTime(timezone string) error {
	command := fmt.Sprintf("ln -sf /usr/share/zoneinfo/%s /etc/localtime", timezone)
	return arch_chroot.Run(command)
}

// Sets up hardware clock to generate /etc/adjtime.
//
// Runs the following command in arch-chroot:
//
//	hwclock --systohc
//
// Can return error types:
//   - PipeError
//   - ArchChrootError
func SetHwClock() error {
	command := "hwclock --systohc"
	return arch_chroot.Run(command)
}

// Checks if the given timezone is a valid.
//
// Can return error types:
//   - TimezoneError
func ValidateTimezone(timezone string) error {
	timezones, err := getAllTimezones()
	if err != nil {
		return TimezoneError{
			Err: err,
		}
	}

	if _, found := slices.BinarySearch(timezones, timezone); !found {
		return TimezoneError{
			Err: errors.New("Invalid timezone"),
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
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return strings.Split(string(stdoutBytes), "\n"), nil
}
