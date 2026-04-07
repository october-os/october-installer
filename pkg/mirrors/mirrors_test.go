package mirrors

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateValidCountry(t *testing.T) {
	originalMirrorListPath := mirrorlistFile
	defer func() {
		mirrorlistFile = originalMirrorListPath
	}()

	createAndSetMirrorListMock(t)

	country := "Canada"

	err := ValidateCountry(country)

	assert.NoError(t, err)
}

func TestValidateInvalidCountry(t *testing.T) {
	originalMirrorListPath := mirrorlistFile
	defer func() {
		mirrorlistFile = originalMirrorListPath
	}()

	createAndSetMirrorListMock(t)

	country := "Montreal"

	err := ValidateCountry(country)

	assert.Error(t, err)
	assert.ErrorAs(t, err, &MirrorListError{}, "Error should be of type mirror list error")
}

func TestValidateAmbiguousString(t *testing.T) {
	originalMirrorListPath := mirrorlistFile
	defer func() {
		mirrorlistFile = originalMirrorListPath
	}()

	createAndSetMirrorListMock(t)

	country := "#Server"

	err := ValidateCountry(country)

	assert.Error(t, err)
	assert.ErrorAs(t, err, &MirrorListError{}, "Error should be of type mirror list error")
}

func createAndSetMirrorListMock(t *testing.T) {
	t.Helper()
	tempDir := t.TempDir()
	tempFile, _ := os.CreateTemp(tempDir, "mirrorlist")

	testdata, _ := os.Open("testdata/mirrorlist")
	defer testdata.Close()

	io.Copy(tempFile, testdata)
	mirrorlistFile = tempFile.Name()
}
