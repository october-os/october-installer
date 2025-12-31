package mirrors

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"
)

const mirrorlistFile string = "/etc/pacman.d/mirrorlist"

func SetMirrorList(countries []string) error {
	mirrorMap, err := getMirrors()
	if err != nil {
		return err
	}

	fullCountrySlice := make([]string, 0, len(mirrorMap))
	for k := range mirrorMap {
		fullCountrySlice = append(fullCountrySlice, k)
	}

	if err := areCountriesValid(countries, fullCountrySlice); err != nil {
		return err
	}

	return saveMirrorlist(countries, mirrorMap)
}

func saveMirrorlist(countries []string, mirrorMap map[string][]string) error {
	file, err := os.Create(mirrorlistFile)
	if err != nil {
		return MirrorListError{
			err: err,
		}
	}
	defer file.Close()

	for _, country := range countries {
		for _, server := range mirrorMap[country] {
			if _, err := file.WriteString(server + "\n"); err != nil {
				return MirrorListError{
					err: err,
				}
			}
		}
	}

	return nil
}

func areCountriesValid(countries []string, fullCountryList []string) error {
	for _, country := range countries {
		if !slices.Contains(fullCountryList, country) {
			return MirrorListError{
				err: fmt.Errorf("country %s not present in mirrors", country),
			}
		}
	}

	return nil
}

func getMirrors() (map[string][]string, error) {
	file, err := os.Open(mirrorlistFile)
	if err != nil {
		return nil, MirrorListError{
			err: err,
		}
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
