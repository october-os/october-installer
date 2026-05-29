package json_parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/october-os/october-installer/pkg/mirrors"
	"github.com/stretchr/testify/assert"
)

func TestValidJson(t *testing.T) {
	originalMirrorList := mirrors.MirrorlistFile
	mirrors.MirrorlistFile = "../mirrors/testdata/mirrorlist"
	defer func() {
		mirrors.MirrorlistFile = originalMirrorList
	}()

	fileData, err := readJsonFile("valid")
	assert.Nil(t, err)

	installation, err := ParseJson(string(fileData))
	assert.Nil(t, err)

	assert.NotEmpty(t, installation)
	assert.Len(t, installation.Drives, 2)
	assert.Equal(t, "/dev/sda", installation.Drives[0].Path)
	assert.False(t, installation.Drives[0].Append)
	assert.Len(t, installation.Drives[0].Partitions, 3)

	assert.Len(t, installation.Users, 2)
	assert.Equal(t, "testuser", installation.Users[0].Username)
	assert.True(t, installation.Users[0].Sudoer)
	assert.Equal(t, "secondtestuser", installation.Users[1].Username)
	assert.False(t, installation.Users[1].Sudoer)

	assert.Len(t, installation.MirrorCountries, 1)
	assert.Equal(t, "Canada", installation.MirrorCountries[0])
	assert.Equal(t, "America/Montreal", installation.Timezone)
	assert.Equal(t, "en_us.UTF-8", installation.Locale)
	assert.Equal(t, "testhostname", installation.Hostname)
	assert.Equal(t, "test", installation.RootPassword)
	assert.False(t, installation.BestEffortGpu)
	assert.Len(t, installation.ExtraPackages.OfficialRepositories, 2)
	assert.Len(t, installation.ExtraPackages.AUR, 1)
}

func readJsonFile(name string) ([]byte, error) {
	filename := fmt.Sprintf("testdata/%s.json", name)
	return os.ReadFile(filename)
}
