package user

import (
	"errors"
	"fmt"
	"testing"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/mocks"
	"github.com/stretchr/testify/assert"
)

// archChrootCommandBuilder wraps a command in the chroot prefix, matching what
// utils.NewCommandExecutor produces when run through CommandExecutorMock.String().
func archChrootCommandBuilder(command string) string {
	return fmt.Sprintf("/usr/bin/arch-chroot /mnt /bin/bash -c %s", command)
}

// ----- Validate tests (no mocking needed) -----

func TestValidateHappy(t *testing.T) {
	u := &User{
		Username: "testuser",
		Password: "P@ssw0rd",
		Homepath: "/home/testuser",
		Sudoer:   true,
	}
	assert.Nil(t, u.Validate())
}

func TestValidateEmptyUsername(t *testing.T) {
	u := &User{
		Username: "",
		Password: "P@ssw0rd",
	}
	err := u.Validate()
	assert.Error(t, err)
	var newErr NewUserError
	assert.ErrorAs(t, err, &newErr)
}

func TestValidateWhitespaceUsername(t *testing.T) {
	u := &User{
		Username: "   ",
		Password: "P@ssw0rd",
	}
	err := u.Validate()
	assert.Error(t, err)
	var newErr NewUserError
	assert.ErrorAs(t, err, &newErr)
}

func TestValidateEmptyPassword(t *testing.T) {
	u := &User{
		Username: "testuser",
		Password: "",
	}
	err := u.Validate()
	assert.Error(t, err)
	var newErr NewUserError
	assert.ErrorAs(t, err, &newErr)
}

func TestValidateRelativeHomepath(t *testing.T) {
	u := &User{
		Username: "testuser",
		Password: "P@ssw0rd",
		Homepath: "home/testuser",
	}
	err := u.Validate()
	assert.Error(t, err)
	var newErr NewUserError
	assert.ErrorAs(t, err, &newErr)
}

func TestValidateNoHomepath(t *testing.T) {
	u := &User{
		Username: "testuser",
		Password: "P@ssw0rd",
		Homepath: "",
	}
	assert.Nil(t, u.Validate())
}

// ----- SetRootPassword tests -----

func TestSetRootPassword(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := SetRootPassword("S3cret")
	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder("echo S3cret | passwd -s"),
		mocks.CommandExecutorGot[0],
	)
}

func TestSetRootPasswordEscapesDollar(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	err := SetRootPassword("P$as$s")
	assert.NoError(t, err)
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder("echo P\\$as\\$s | passwd -s"),
		mocks.CommandExecutorGot[0],
	)
}

func TestSetRootPasswordError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	err := SetRootPassword("S3cret")
	assert.Error(t, err)
	var newErr NewUserError
	assert.ErrorAs(t, err, &newErr)
	assert.Len(t, mocks.CommandExecutorGot, 1)
}

// ----- CreateUser tests -----

func TestCreateUserNonSudoer(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	u := &User{
		Username: "testuser",
		Password: "S3cret",
		Homepath: "/home/testuser",
		Sudoer:   false,
	}

	assert.NoError(t, CreateUser(u))
	assert.Len(t, mocks.CommandExecutorGot, 2)
	assert.Equal(
		t,
		archChrootCommandBuilder("useradd -m testuser -d /home/testuser"),
		mocks.CommandExecutorGot[0],
	)
	assert.Equal(
		t,
		archChrootCommandBuilder("echo S3cret | passwd testuser -s"),
		mocks.CommandExecutorGot[1],
	)
}

func TestCreateUserSudoer(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	u := &User{
		Username: "testuser",
		Password: "S3cret",
		Homepath: "/home/testuser",
		Sudoer:   true,
	}

	assert.NoError(t, CreateUser(u))
	assert.Len(t, mocks.CommandExecutorGot, 3)
	assert.Equal(
		t,
		archChrootCommandBuilder("useradd -m testuser -d /home/testuser"),
		mocks.CommandExecutorGot[0],
	)
	assert.Equal(
		t,
		archChrootCommandBuilder("echo S3cret | passwd testuser -s"),
		mocks.CommandExecutorGot[1],
	)
	assert.Equal(
		t,
		archChrootCommandBuilder("usermod -aG wheel testuser"),
		mocks.CommandExecutorGot[2],
	)
}

func TestCreateUserDefaultHomepath(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	u := &User{
		Username: "testuser",
		Password: "S3cret",
		Homepath: "",
		Sudoer:   false,
	}

	assert.NoError(t, CreateUser(u))
	assert.Len(t, mocks.CommandExecutorGot, 2)
	assert.Equal(
		t,
		archChrootCommandBuilder("useradd -m testuser -d /home/testuser"),
		mocks.CommandExecutorGot[0],
	)
	assert.Equal(
		t,
		"/home/testuser",
		u.Homepath,
		"Homepath should be filled in when empty",
	)
}

func TestCreateUserEscapesPassword(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	u := &User{
		Username: "testuser",
		Password: "P@$w@rd$",
		Homepath: "/home/testuser",
		Sudoer:   false,
	}

	assert.NoError(t, CreateUser(u))
	assert.Len(t, mocks.CommandExecutorGot, 2)
	assert.Equal(
		t,
		archChrootCommandBuilder("echo P@\\$w@rd\\$ | passwd testuser -s"),
		mocks.CommandExecutorGot[1],
	)
}

func TestCreateUserError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	u := &User{
		Username: "testuser",
		Password: "S3cret",
		Homepath: "/home/testuser",
		Sudoer:   false,
	}

	err := CreateUser(u)
	assert.Error(t, err)
	var newErr NewUserError
	assert.ErrorAs(t, err, &newErr)
	assert.Len(t, mocks.CommandExecutorGot, 1)
}

// ----- Private helper tests -----

func TestUserAdd(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	u := &User{
		Username: "testuser",
		Homepath: "/home/testuser",
	}

	assert.NoError(t, userAdd(u))
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder("useradd -m testuser -d /home/testuser"),
		mocks.CommandExecutorGot[0],
	)
}

func TestUserAddDefaultsHomepath(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	u := &User{
		Username: "testuser",
		Homepath: "",
	}

	assert.NoError(t, userAdd(u))
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder("useradd -m testuser -d /home/testuser"),
		mocks.CommandExecutorGot[0],
	)
	assert.Equal(t, "/home/testuser", u.Homepath)
}

func TestUserAddError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	u := &User{
		Username: "testuser",
		Homepath: "/home/testuser",
	}

	assert.Error(t, userAdd(u))
	assert.Len(t, mocks.CommandExecutorGot, 1)
}

func TestSetPassword(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	assert.NoError(t, setPassword("testuser", "S3cret"))
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder("echo S3cret | passwd testuser -s"),
		mocks.CommandExecutorGot[0],
	)
}

func TestSetPasswordEscapesDollar(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	testCases := []struct {
		password string
		expected string
	}{
		{"a$b", "a\\$b"},
		{"$$", "\\$\\$"},
		{"$6$salt$hash", "\\$6\\$salt\\$hash"},
		{"no-dollars", "no-dollars"},
	}

	for _, tc := range testCases {
		mocks.CommandExecutorGot = []string{}
		assert.NoError(t, setPassword("testuser", tc.password))
		assert.Len(t, mocks.CommandExecutorGot, 1)
		assert.Equal(
			t,
			archChrootCommandBuilder(fmt.Sprintf("echo %s | passwd testuser -s", tc.expected)),
			mocks.CommandExecutorGot[0],
		)
	}
}

func TestSetPasswordError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	assert.Error(t, setPassword("testuser", "S3cret"))
	assert.Len(t, mocks.CommandExecutorGot, 1)
}

func TestAddToSudoer(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
	}()

	assert.NoError(t, addToSudoer("testuser"))
	assert.Len(t, mocks.CommandExecutorGot, 1)
	assert.Equal(
		t,
		archChrootCommandBuilder("usermod -aG wheel testuser"),
		mocks.CommandExecutorGot[0],
	)
}

func TestAddToSudoerError(t *testing.T) {
	originalExecutor := arch_chroot.CommandExecutor
	arch_chroot.CommandExecutor = mocks.NewCommandExecutorMock
	mocks.ReturnError = true
	defer func() {
		arch_chroot.CommandExecutor = originalExecutor
		mocks.CommandExecutorGot = []string{}
		mocks.ReturnError = false
	}()

	assert.Error(t, addToSudoer("testuser"))
	assert.Len(t, mocks.CommandExecutorGot, 1)
}

// ----- NewUserError tests -----

func TestNewUserErrorError(t *testing.T) {
	wrapped := errors.New("some internal failure")
	e := NewUserError{err: wrapped}
	assert.Equal(t, "New user error: error=some internal failure", e.Error())
}

func TestNewUserErrorUnwrap(t *testing.T) {
	wrapped := errors.New("some internal failure")
	e := NewUserError{err: wrapped}
	assert.True(t, errors.Is(e.Unwrap(), wrapped))
	assert.True(t, errors.Is(e, wrapped))
}
