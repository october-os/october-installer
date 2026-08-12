package mirrors

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValidCountry(t *testing.T) {
	originalMirrorListPath := MirrorlistFile
	defer func() {
		MirrorlistFile = originalMirrorListPath
		mirrorMap = nil
	}()

	createAndSetMirrorListMock(t)

	country := "Canada"

	err := ValidateCountry(country)

	assert.NoError(t, err)
}

func TestValidateInvalidCountry(t *testing.T) {
	originalMirrorListPath := MirrorlistFile
	defer func() {
		MirrorlistFile = originalMirrorListPath
		mirrorMap = nil
	}()

	createAndSetMirrorListMock(t)

	country := "Montreal"

	err := ValidateCountry(country)

	assert.Error(t, err)
	assert.ErrorAs(t, err, &MirrorListError{}, "Error should be of type mirror list error")
}

func TestValidateAmbiguousString(t *testing.T) {
	originalMirrorListPath := MirrorlistFile
	defer func() {
		MirrorlistFile = originalMirrorListPath
		mirrorMap = nil
	}()

	createAndSetMirrorListMock(t)

	country := "#Server"

	err := ValidateCountry(country)

	assert.Error(t, err)
	assert.ErrorAs(t, err, &MirrorListError{}, "Error should be of type mirror list error")
}

func TestSetMirrorValidList(t *testing.T) {
	originalMirrorListPath := MirrorlistFile
	defer func() {
		MirrorlistFile = originalMirrorListPath
		mirrorMap = nil
	}()

	createAndSetMirrorListMock(t)

	country := []string{"France", "Canada"}

	err := SetMirrorList(country)

	assert.NoError(t, err)
	assert.NotNil(t, mirrorMap[country[0]], "France mirrors should be in the mirror map")
	assert.NotNil(t, mirrorMap[country[1]], "Canada mirrors should be in the mirror map")

	content, _ := os.ReadFile(MirrorlistFile)
	assert.NotNil(t, content)

	for _, v := range country {
		for _, server := range mirrorMap[v] {
			assert.Contains(t, string(content), server)
		}
	}
}

func createAndSetMirrorListMock(t *testing.T) {
	t.Helper()
	tempDir := t.TempDir()
	tempFile, _ := os.CreateTemp(tempDir, "mirrorlist")

	testdata, _ := os.Open("testdata/mirrorlist")
	defer testdata.Close()

	io.Copy(tempFile, testdata)
	MirrorlistFile = tempFile.Name()
}
