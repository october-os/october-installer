package keyring

import (
	"testing"

	"github.com/october-os/october-installer/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdatingKeyring(t *testing.T) {
	commandExecutor = utils.NewCommandExecutorMock
	defer func() {
		commandExecutor = utils.NewCommandExecutor
		utils.CommandExecutorGot = []string{}
	}()

	expectedCommand := "/usr/bin/pacman -Sy --noconfirm archlinux-keyring"

	err := UpdateKeyRing()

	assert.NoError(t, err)
	assert.Len(t, utils.CommandExecutorGot, 1, "No command ran during the test")
	assert.Equal(t, expectedCommand, utils.CommandExecutorGot[0])
}

func TestUpdatingKeyringWithError(t *testing.T) {
	commandExecutor = utils.NewCommandExecutorMock
	utils.ReturnMockError = true
	defer func() {
		commandExecutor = utils.NewCommandExecutor
		utils.ReturnMockError = false
	}()

	err := UpdateKeyRing()

	assert.Error(t, err)
	assert.ErrorAs(t, err, &KeyringUpdateError{}, "Error should be of type KeyringUpdateError")
}
