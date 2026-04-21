package fstab

import (
	"os"
	"testing"

	"github.com/october-os/october-installer/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFstab(t *testing.T) {
	want := "/usr/bin/genfstab -U /mnt"
	fileContentWant := "testing"

	tempDir := t.TempDir()
	tempFile, _ := os.CreateTemp(tempDir, "test")

	mocks.OutputReturn = fileContentWant
	normalCommandExecutor := commandExecutor
	normalFstabFilePath := fstabPath
	defer func() {
		commandExecutor = normalCommandExecutor
		fstabPath = normalFstabFilePath
		mocks.CommandExecutorGot = []string{}
		mocks.OutputReturn = ""
	}()

	fstabPath = tempFile.Name()
	commandExecutor = mocks.NewCommandExecutorMock

	err := GenerateFstab()

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1, "no command ran during this test.")
	assert.Equal(t, want, mocks.CommandExecutorGot[0])

	got, _ := os.ReadFile(tempFile.Name())
	assert.Equal(t, fileContentWant, string(got), "file content doesn't match mock returned bytes array")
}

func TestGenerateFstabError(t *testing.T) {
	normalCommandExecutor := commandExecutor
	defer func() {
		commandExecutor = normalCommandExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	mocks.ReturnError = true
	commandExecutor = mocks.NewCommandExecutorMock

	err := GenerateFstab()

	assert.Error(t, err)
	assert.ErrorAs(t, err, &FstabError{})
}
