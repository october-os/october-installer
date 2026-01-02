package locale

import (
	"fmt"
	"io"
	"os/exec"
	"strings"

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
func ValidateLocale(locale string) (bool, error) {
	command := fmt.Sprintf("cat %s | grep %s", filepath, locale)
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return false, LocaleGenError{
			Err: err,
		}
	}

	if err := cmd.Start(); err != nil {
		return false, LocaleGenError{
			Err: err,
		}
	}

	stdoutBytes, err := io.ReadAll(stdout)
	if err != nil {
		return false, LocaleGenError{
			Err: err,
		}
	}

	if err := cmd.Wait(); err != nil {
		return false, LocaleGenError{
			Err: err,
		}
	}

	return strings.Contains(string(stdoutBytes), locale+" UTF-8"), nil
}
