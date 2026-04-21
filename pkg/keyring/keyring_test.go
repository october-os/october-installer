package keyring

import (
	"testing"

	"github.com/october-os/october-installer/pkg/mocks"
	"github.com/october-os/october-installer/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestUpdatingKeyring(t *testing.T) {
	commandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		commandExecutor = utils.NewCommandExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	expectedCommand := "/usr/bin/pacman -Sy --noconfirm archlinux-keyring"

	err := UpdateKeyRing()

	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1, "No command ran during the test")
	assert.Equal(t, expectedCommand, mocks.CommandExecutorGot[0])
}

func TestUpdatingKeyringWithError(t *testing.T) {
	commandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		commandExecutor = utils.NewCommandExecutor
		mocks.ReturnError = false
	}()

	err := UpdateKeyRing()

	assert.Error(t, err)
	assert.ErrorAs(t, err, &KeyringUpdateError{}, "Error should be of type KeyringUpdateError")
}
