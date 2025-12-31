package mirrors

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

// Absolute path to the mirrorlist file.
const mirrorlistFile string = "/etc/pacman.d/mirrorlist"

// Sets the mirrorlist file with only the servers for the
// given countries and removes all the unused ones.
//
// Can return error types:
//   - MirrorListError
func SetMirrorList(countries []string) error {
	mirrorMap, err := getMirrors()
	if err != nil {
		return MirrorListError{
			err: err,
		}
	}

	fullCountrySlice := make([]string, 0, len(mirrorMap))
	for k := range mirrorMap {
		fullCountrySlice = append(fullCountrySlice, k)
	}

	if err := areCountriesValid(countries, fullCountrySlice); err != nil {
		return MirrorListError{
			err: err,
		}
	}

	if err := saveMirrorlist(countries, mirrorMap); err != nil {
		return MirrorListError{
			err: err,
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

// Checks if the given countries are inside the given full country list.
func areCountriesValid(countries []string, fullCountryList []string) error {
	for _, country := range countries {
		if !slices.Contains(fullCountryList, country) {
			return fmt.Errorf("country %s not present in mirrors", country)
		}
	}

	return nil
}

// Reads the mirrorlist file and returns a map
// that has the country name as the key and a slice of
// all the servers as the value.
func getMirrors() (map[string][]string, error) {
	file, err := os.Open(mirrorlistFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var countryMap map[string][]string = make(map[string][]string)
	var lastCountry string = ""

	for scanner.Scan() {
		line := scanner.Text()
		if line == " " {
			continue
		}

		if country, found := strings.CutPrefix(line, "## "); found {
			lastCountry = country
			countryMap[country] = make([]string, 0)
		} else {
			countryMap[lastCountry] = append(countryMap[lastCountry], strings.TrimPrefix(line, "#"))
		}
	}

	return countryMap, nil
}
