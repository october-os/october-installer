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
func SetRootPassword(password string) error {
	command := fmt.Sprintf("echo %s | passwd -s", password)
	if err := arch_chroot.Run(command); err != nil {
		return err
	}

	return nil
}

// Takes in a user then creates it in the newly installed system.
//
// Errors that can be returned:
//   - PipeError
//   - ArchChrootError
func CreateUser(user *User) error {
	err := userAdd(user.Username, user.Homepath)
	if err != nil {
		return err
	}

	err = setPassword(user.Username, user.Password)
	if err != nil {
		return err
	}

	if user.Sudoer {
		err = addToSudoer(user.Username)
		if err != nil {
			return err
		}
	}

	return nil
}

// Adds the user with the given username to the wheel group to make
// it a sudoer inside the newly installed system.
//
// Errors that can be returned:
//   - PipeError
//   - ArchChrootError
func addToSudoer(username string) error {
	addToWheel := fmt.Sprintf("usermod -aG wheel %s", username)

	err := arch_chroot.Run(addToWheel)
	if err != nil {
		return err
	}

	return nil
}

// Runs useradd with the given username and homepath inside the newly
// installed system.
//
// Errors that can be returned:
//   - PipeError
//   - ArchChrootError
func userAdd(username, homepath string) error {
	createCommand := fmt.Sprintf("useradd -m %s -d %s", username, homepath)

	err := arch_chroot.Run(createCommand)
	if err != nil {
		return err
	}

	return nil
}

// Sets the given user password for the given password inside
// the newly installed system.
//
// Errors that can be returned:
//   - PipeError
//   - ArchChrootError
func setPassword(username, password string) error {
	command := fmt.Sprintf("echo %s | passwd %s -s", password, username)

	err := arch_chroot.Run(command)
	if err != nil {
		return err
	}

	return nil
}
