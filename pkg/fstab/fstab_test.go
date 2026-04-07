package fstab

import (
	"os"
	"testing"

	"github.com/october-os/october-installer/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestGenerateFstab(t *testing.T) {
	want := "/usr/bin/genfstab -U /mnt"
	fileContentWant := "testing"

	tempDir := t.TempDir()
	tempFile, _ := os.CreateTemp(tempDir, "test")

	utils.OutputMockReturn = fileContentWant
	normalCommandExecutor := commandExecutor
	normalFstabFilePath := fstabPath
	defer func() {
		commandExecutor = normalCommandExecutor
		fstabPath = normalFstabFilePath
		utils.CommandExecutorGot = []string{}
		utils.OutputMockReturn = ""
	}()

	fstabPath = tempFile.Name()
	commandExecutor = utils.NewCommandExecutorMock

	err := GenerateFstab()

	assert.NoError(t, err)
	assert.Len(t, utils.CommandExecutorGot, 1, "no command ran during this test.")
	assert.Equal(t, want, utils.CommandExecutorGot[0])

	got, _ := os.ReadFile(tempFile.Name())
	assert.Equal(t, fileContentWant, string(got), "file content doesn't match mock returned bytes array")
}

func TestGenerateFstabError(t *testing.T) {
	normalCommandExecutor := commandExecutor
	defer func() {
		commandExecutor = normalCommandExecutor
		utils.CommandExecutorGot = []string{}
		utils.ReturnMockError = false
	}()

	utils.ReturnMockError = true
	commandExecutor = utils.NewCommandExecutorMock

	err := GenerateFstab()

	assert.Error(t, err)
	assert.ErrorAs(t, err, &FstabError{})
}
