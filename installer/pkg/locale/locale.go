package locale

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/arch-couple/arch-couple-installer/pkg/arch_chroot"
)

const FILEPATH string = "/mnt/etc/locale.gen"

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

	if err := arch_chroot.Run("locale-gen"); err != nil {
		return err
	}

	return nil
}

func doesLocaleExist(locale string, allLocales []string) (bool, int) {
	for i, line := range allLocales {
		if strings.Contains(line, locale) {
			return true, i
		}
	}

	return false, 0
}

func saveLocaleGen(content []string) error {
	if err := arch_chroot.Run("rm /etc/locale.gen"); err != nil {
		return err
	}

	for _, line := range content {
		fmt.Println(line)
		command := fmt.Sprintf("echo \"%s\" >> /etc/locale.gen", line)
		if err := arch_chroot.Run(command); err != nil {
			return err
		}
	}

	return nil
}

func loadLocaleGen() ([]string, error) {
	f, err := os.Open(FILEPATH)
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
