// Package user provides the struct representing a user that
// needs to be created and the functions to create it in the newly
// installed system.
package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

// User represents a user that needs to be created.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Homepath string `json:"homepath"`
	Sudoer   bool   `json:"sudoer"`
}

// Validates if the user is a valid one or if it contains values that
// aren't valid.
func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" || strings.TrimSpace(u.Password) == "" {
		return NewUserError{
			err: errors.New("Can't create user with empty username or password"),
		}
	}

	if strings.TrimSpace(u.Homepath) == "" {
		u.Homepath = fmt.Sprintf("/home/%s", u.Username)
	} else if !strings.HasPrefix(u.Homepath, "/") {
		return NewUserError{
			err: errors.New("Provide a valid directory for user home path"),
		}
	}

	return nil
}

// Sets the given password for the root user.
// Errors that can be returned:
//   - NewUserError
func SetRootPassword(password string) error {
	command := fmt.Sprintf("echo %s | passwd -s", password)
	if err := arch_chroot.Run(command); err != nil {
		return NewUserError{err: err}
	}

	return nil
}

// Takes in a user then creates it in the newly installed system.
//
// Errors that can be returned:
//   - NewUserError
func CreateUser(user *User) error {
	if err := userAdd(user.Username, user.Homepath); err != nil {
		return NewUserError{err: err}
	}

	if err := setPassword(user.Username, user.Password); err != nil {
		return NewUserError{err: err}
	}

	if user.Sudoer {

		if err := addToSudoer(user.Username); err != nil {
			return NewUserError{err: err}
		}
	}

	return nil
}

// Adds the user with the given username to the wheel group to make
// it a sudoer inside the newly installed system.
func addToSudoer(username string) error {
	addToWheel := fmt.Sprintf("usermod -aG wheel %s", username)
	return arch_chroot.Run(addToWheel)
}

// Runs useradd with the given username and homepath inside the newly
// installed system.
func userAdd(username, homepath string) error {
	createCommand := fmt.Sprintf("useradd -m %s -d %s", username, homepath)
	return arch_chroot.Run(createCommand)
}

// Sets the given user password for the given password inside
// the newly installed system.
func setPassword(username, password string) error {
	command := fmt.Sprintf("echo %s | passwd %s -s", password, username)
	return arch_chroot.Run(command)
}
