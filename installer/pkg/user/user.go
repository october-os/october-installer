// Package user provides the struct representing a user that
// needs to be created and the functions to create it in the newly
// installed system.
package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arch-couple/arch-couple-installer/pkg/arch_chroot"
)

// User represents a user that needs to be created.
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Homepath string `json:"homepath"`
	Sudoer   bool   `json:"sudoer"`
}

// Constructor for User struct.
//
// If homepath is empty, it will be set to: /home/[username]
//
// Will either return:
//   - *User
//   - NewUserError: When an incorrect value is passed
func New(username, password, homepath string, sudoer bool) (*User, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, NewUserError{
			err: errors.New("Can't create user with empty username or password"),
		}
	}

	if strings.TrimSpace(homepath) == "" {
		homepath = fmt.Sprintf("/home/%s", username)
	} else if strings.HasPrefix(homepath, "/") {
		return nil, NewUserError{
			err: errors.New("Provide a valid directory for user home path"),
		}
	}

	return &User{
		Username: username,
		Password: password,
		Homepath: homepath,
		Sudoer:   sudoer,
	}, nil
}

// Takes in a user then creates it in the newly installed system.
//
// Errors that can be returned:
//   - PipeError
//   - ArchChrootError
func CreateUser(user User) error {
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

// Adds the wheel user group as system admins in /etc/sudoers
// inside the newly installed system.
//
// Errors that can be returned:
//   - PipeError
//   - ArchChrootError
func SetupSudoerFile() error {
	wheelLine := "%wheel      ALL=(ALL:ALL) ALL"
	command := fmt.Sprintf("echo \"%s\" >> /etc/sudoers", wheelLine)

	err := arch_chroot.Run(command)
	if err != nil {
		return err
	}

	return nil
}

// Adds the user with the given username to the wheel group to make
// it a sudoer inside the newly installed system. Make sure to run
// SetupSudoerFile() before running this.
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
