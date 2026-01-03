package json_parser

import (
	"encoding/json"
	"fmt"

	"github.com/october-os/october-installer/pkg/partition"
	"github.com/october-os/october-installer/pkg/user"
)

type Installation struct {
	Drives       []partition.Drive `json:"drives"`
	Users        []user.User       `json:"users"`
	Mirrors      []string          `json:"mirrors"`
	Timezone     string            `json:"timezone"`
	Locale       string            `json:"locale"`
	Hostname     string            `json:"hostname"`
	RootPassword string            `json:"rootPassword"`
}

func ParseJson(jsonString string) (*Installation, error) {
	jsonBytes := []byte(jsonString)
	var installation Installation
	if err := json.Unmarshal(jsonBytes, &installation); err != nil {
		return nil, &JsonParsingError{
			Err: fmt.Errorf("error parsing json: error=%s", err.Error()),
		}
	}
	for _, drive := range installation.Drives {
		if err := drive.Validate(); err != nil {
			return nil, err
		}
	}
	for _, user := range installation.Users {
		if err := user.Validate(); err != nil {
			return nil, err
		}
	}
	// TODO: validate mirrors, timezone, locale, hostname
	return &installation, nil
}
