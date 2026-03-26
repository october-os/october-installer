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

	assert.Nil(t, err)
	assert.Len(t, utils.CommandExecutorGot, 1, "no command ran during this test.")
	assert.Equal(t, want, utils.CommandExecutorGot[0])

	actualFileContent := make([]byte, len([]byte(fileContentWant)))
	_, err = tempFile.Read(actualFileContent)
	assert.Nil(t, err)
	assert.Equal(t, []byte(fileContentWant), actualFileContent, "file content doesn't match mock returned bytes array")
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

	assert.NotNil(t, err)
	assert.ErrorAs(t, err, &FstabError{})
}
