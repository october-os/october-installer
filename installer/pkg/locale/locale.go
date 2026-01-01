package locale

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// Absolute file path to locale.gen
const filepath string = "/mnt/etc/locale.gen"

// Checks if the given locale exist, then uncomment its and
// sets up the locales.
//
// The given locale must be in UTF-8 and in the same
// format as it is inside /etc/locale.gen before the space.
//
// Can return error types:
//   - LocaleGenError
//   - PipeError
//   - ArchChrootError
func GenerateLocales(locale string) error {
	locales, err := loadLocaleGen()
	if err != nil {
		return LocaleGenError{
			Err: err,
		}
	}

	if exists, index := doesLocaleExist(locale, locales); exists {
		locales[index] = strings.TrimLeft(locales[index], "#")
		splitedLocale := strings.Split(locales[index], " ")

		command := fmt.Sprintf("echo LANG=%s > /etc/locale.conf", splitedLocale[0])
		if err := arch_chroot.Run(command); err != nil {
			return err
		}
	} else {
		return LocaleGenError{
			Err: errors.New("Invalid locale"),
		}
	}

	if err := saveLocaleGen(locales); err != nil {
		return LocaleGenError{
			Err: err,
		}
	}

	return arch_chroot.Run("locale-gen")
}

// Checks if the given UTF-8 locale exist insides /etc/locale.gen.
//
// It returns if the locale was found and its index
// inside the array.
func doesLocaleExist(locale string, allLocales []string) (bool, int) {
	for i, line := range allLocales {
		splittedString := strings.Split(line, " ")
		if splittedString[0] == "#"+locale && splittedString[1] == "UTF-8" {
			return true, i
		}
	}

	return false, 0
}

// Saves the updated content inside /etc/locale.gen using
// arch-chroot.
//
// Updates the file by doing:
//
//	cat > /etc/locale.gen << EOF\n[content]EOF
func saveLocaleGen(content []string) error {
	var contentStr strings.Builder
	for _, line := range content {
		contentStr.WriteString(line)
		contentStr.WriteString("\n")
	}

	command := fmt.Sprintf("cat > /etc/locale.gen << EOF\n%sEOF", contentStr.String())

	return arch_chroot.Run(command)
}

// Reads the content of /etc/locale.gen then
// returns an array of all the lines inside the files.
func loadLocaleGen() ([]string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer f.Close()
	scanner := bufio.NewScanner(f)
	var fileContent []string

	for scanner.Scan() {
		line := scanner.Text()
		fileContent = append(fileContent, line)
	}

	return fileContent, nil
}
