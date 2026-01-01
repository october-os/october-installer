package timezone

import (
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Checks if the timezone exist then sets it up inside
// the new install.
//
// Can return error types:
//   - TimezoneError
//   - PipeError
//   - ArchChrootError
func SetTime(timezone string) error {
	isValid, err := isTimezoneValid(timezone)
	if err != nil {
		return TimezoneError{
			Err: err,
		}
	} else if !isValid {
		return TimezoneError{
			Err: errors.New("Timezone isn't valid"),
		}
	}

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

// Checks if the given timezone is a valid one with a
// binary search accross all timezones.
func isTimezoneValid(timezone string) (bool, error) {
	timezones, err := getAllTimezones()
	if err != nil {
		return false, err
	}

	low := 0
	high := len(timezones) - 1

	for low <= high {
		mid := low + (high-low)/2

		if timezones[mid] == timezone {
			return true, nil
		}

		if timezones[mid] < timezone {
			low = mid + 1
		} else {
			high = mid - 1
		}
	}

	return false, nil
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
