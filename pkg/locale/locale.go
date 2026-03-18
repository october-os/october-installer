package locale

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Absolute file path to locale.gen
const filepath string = "/etc/locale.gen"

// Uncomment and sets up the locales.
//
// The given locale must be in UTF-8 and in the same
// format as it is inside /etc/locale.gen before the space.
//
// Can return error types:
//   - LocaleError
func GenerateLocales(locale string) error {
	sedCmd := fmt.Sprintf("sed -i 's/#%s UTF-8/%s UTF-8/' %s", locale, locale, filepath)
	localeConfCmd := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", locale)
	localegenCmd := "locale-gen"

	if err := arch_chroot.Run(
		sedCmd,
		localeConfCmd,
		localegenCmd); err != nil {
		return LocaleError{err: err}
	}
	return nil
}

// Checks if the given UTF-8 locale exist insides /etc/locale.gen.
//
// Can return error types:
//   - LocaleError
func ValidateLocale(locale string) error {
	command := fmt.Sprintf("cat %s | grep -w \"#%s UTF-8\"", filepath, locale)
	cmd := exec.Command("/bin/bash", "-c", command)

	if err := cmd.Run(); err != nil {
		if cmd.ProcessState.ExitCode() == 1 { // not found
			return LocaleError{
				err: errors.New("Invalid locale"),
			}
		}

		return LocaleError{
			err: err,
		}
	}

	return nil
}
