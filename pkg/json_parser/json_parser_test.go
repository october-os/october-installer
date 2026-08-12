package json_parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/october-os/october-installer/pkg/hostname"
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
	assert.Equal(t, 1, installation.Drives[0].Partitions[0].Size.Amount)
	assert.Equal(t, "GiB", installation.Drives[0].Partitions[0].Size.Unit)
	assert.Equal(t, "C12A7328-F81F-11D2-BA4B-00A0C93EC93B", installation.Drives[0].Partitions[0].PartitionType)
	assert.Equal(t, 4, installation.Drives[0].Partitions[1].Size.Amount)
	assert.Equal(t, "GiB", installation.Drives[0].Partitions[1].Size.Unit)
	assert.Equal(t, "0657FD6D-A4AB-43C4-84E5-0933C84B4F4F", installation.Drives[0].Partitions[1].PartitionType)
	assert.True(t, installation.Drives[0].Partitions[2].Size.TakeRemaining)
	assert.Equal(t, "ext4", installation.Drives[0].Partitions[2].FileSystem)
	assert.Equal(t, "4F68BCE3-E8CD-4DB1-96E7-FBCAF984B709", installation.Drives[0].Partitions[2].PartitionType)

	assert.Equal(t, "/dev/sdb", installation.Drives[1].Path)
	assert.True(t, installation.Drives[1].Append)
	assert.Len(t, installation.Drives[1].Partitions, 2)
	assert.Equal(t, 3, installation.Drives[1].Partitions[0].Size.Amount)
	assert.Equal(t, "GiB", installation.Drives[1].Partitions[0].Size.Unit)
	assert.Equal(t, "0FC63DAF-8483-4772-8E79-3D69D8477DE4", installation.Drives[1].Partitions[0].PartitionType)
	assert.Equal(t, "ext4", installation.Drives[1].Partitions[0].FileSystem)
	assert.Equal(t, "/data", installation.Drives[1].Partitions[0].MountPoint)
	assert.Equal(t, 5, installation.Drives[1].Partitions[1].Size.Amount)
	assert.Equal(t, "GiB", installation.Drives[1].Partitions[1].Size.Unit)
	assert.Equal(t, "933AC7E1-2EB4-4F13-B844-0E14E2AEF915", installation.Drives[1].Partitions[1].PartitionType)
	assert.Equal(t, "btrfs", installation.Drives[1].Partitions[1].FileSystem)
	assert.Equal(t, "/home/testuser", installation.Drives[1].Partitions[1].MountPoint)

	assert.Len(t, installation.Users, 2)
	assert.Equal(t, "testuser", installation.Users[0].Username)
	assert.Equal(t, "test", installation.Users[0].Password)
	assert.True(t, installation.Users[0].Sudoer)
	assert.Equal(t, "secondtestuser", installation.Users[1].Username)
	assert.False(t, installation.Users[1].Sudoer)
	assert.Equal(t, "test", installation.Users[1].Password)

	assert.Len(t, installation.MirrorCountries, 1)
	assert.Equal(t, "Canada", installation.MirrorCountries[0])
	assert.Equal(t, "America/Montreal", installation.Timezone)
	assert.Equal(t, "en_US.UTF-8", installation.Locale)
	assert.Equal(t, "testhostname", installation.Hostname)
	assert.Equal(t, "test", installation.RootPassword)
	assert.False(t, installation.BestEffortGpu)
	assert.Len(t, installation.ExtraPackages.OfficialRepositories, 2)
	assert.Equal(t, "cowsay", installation.ExtraPackages.OfficialRepositories[0])
	assert.Equal(t, "sl", installation.ExtraPackages.OfficialRepositories[1])
	assert.Len(t, installation.ExtraPackages.AUR, 1)
	assert.Equal(t, "neofetch", installation.ExtraPackages.AUR[0])
}

func TestMalformedJson(t *testing.T) {
	originalMirrorList := mirrors.MirrorlistFile
	mirrors.MirrorlistFile = "../mirrors/testdata/mirrorlist"
	defer func() {
		mirrors.MirrorlistFile = originalMirrorList
	}()

	fileData, err := readJsonFile("malformed_no_end_brace")
	assert.Nil(t, err)

	installation, err := ParseJson(string(fileData))
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &JsonParsingError{})

	assert.Nil(t, installation)
}

func TestJsonWithEmptyHostname(t *testing.T) {
	originalMirrorList := mirrors.MirrorlistFile
	mirrors.MirrorlistFile = "../mirrors/testdata/mirrorlist"
	defer func() {
		mirrors.MirrorlistFile = originalMirrorList
	}()

	fileData, err := readJsonFile("emptyHostname")
	assert.Nil(t, err)

	installation, err := ParseJson(string(fileData))
	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &JsonParsingError{})
	assert.ErrorAs(t, err, &hostname.HostnameError{})

	assert.Nil(t, installation)
}

func readJsonFile(name string) ([]byte, error) {
	filename := fmt.Sprintf("testdata/%s.json", name)
	return os.ReadFile(filename)
}
