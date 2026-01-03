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
//   - PipeError
//   - ArchChrootError
func GenerateLocales(locale string) error {
	sedCmd := fmt.Sprintf("sed -i 's/#%s UTF-8/%s UTF-8/' %s", locale, locale, filepath)
	localeConfCmd := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", locale)
	localegenCmd := "locale-gen"

	command := fmt.Sprintf("%s && %s && %s", sedCmd, localeConfCmd, localegenCmd)
	return arch_chroot.Run(command)
}

// Checks if the given UTF-8 locale exist insides /etc/locale.gen.
//
// Can return error types:
//   - LocaleGenError
func ValidateLocale(locale string) error {
	command := fmt.Sprintf("cat %s | grep \"%s UTF-8\"", filepath, locale)
	cmd := exec.Command("/bin/bash", "-c", command)

	if err := cmd.Run(); err != nil {
		if cmd.ProcessState.ExitCode() == 1 { // not found
			return LocaleGenError{
				Err: errors.New("Invalid locale"),
			}
		} else {
			return LocaleGenError{
				Err: err,
			}
		}
	}

	return nil
}
