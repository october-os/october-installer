package locale

import (
	"bytes"
	"fmt"
	"os"

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
	content, err := os.ReadFile(filepath)
	if err != nil {
		return LocaleError{err: err}
	}

	if !bytes.Contains(content, []byte(fmt.Sprintf("#%s UTF-8", locale))) {
		return LocaleError{
			err: fmt.Errorf("%s is not a valid locale file", locale),
		}
	}

	return nil
}
