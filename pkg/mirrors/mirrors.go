package mirrors

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Absolute path to the mirrorlist file.
var mirrorlistFile string = "/etc/pacman.d/mirrorlist"

// mirrorMap Map of all the mirrors in the mirrorlist
// file. Key is the country and the value is an array
// of servers.
var mirrorMap map[string][]string = nil

// Sets the mirrorlist file with only the servers for the
// given countries and removes all the unused ones.
//
// Can return error types:
//   - MirrorListError
func SetMirrorList(countries []string) error {
	if mirrorMap == nil {
		err := getMirrors()
		if err != nil {
			return MirrorListError{
				err: err,
			}
		}
	}

	if err := saveMirrorlist(countries, mirrorMap); err != nil {
		return MirrorListError{
			err: err,
		}
	}

	return nil
}

// Checks if the given country is present inside mirrorlist.
//
// Can return errors of types:
//   - MirrorListError
func ValidateCountry(country string) error {
	if mirrorMap == nil {
		err := getMirrors()
		if err != nil {
			return MirrorListError{err: err}
		}
	}

	_, found := mirrorMap[country]
	if !found {
		return MirrorListError{
			err: fmt.Errorf("%s is not a valid country", country),
		}
	}

	return nil
}

// Saves all the servers of the given countries inside the
// mirrorlist file.
func saveMirrorlist(countries []string, mirrorMap map[string][]string) error {
	file, err := os.Create(mirrorlistFile)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, country := range countries {
		for _, server := range mirrorMap[country] {
			if _, err := file.WriteString(server + "\n"); err != nil {
				return err
			}
		}
	}

	return nil
}

// getMirrors reads the mirrorlist file and populate mirrorMap.
func getMirrors() error {
	file, err := os.Open(mirrorlistFile)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	mirrorMap = make(map[string][]string)
	var lastCountry string = ""

	for scanner.Scan() {
		line := scanner.Text()

		if country, found := strings.CutPrefix(line, "## "); found {
			lastCountry = country
			mirrorMap[country] = make([]string, 0)
		} else {
			mirrorMap[lastCountry] = append(mirrorMap[lastCountry], strings.TrimPrefix(line, "#"))
		}
	}

	return nil
}
