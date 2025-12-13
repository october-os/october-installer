package user

import (
	"errors"
	"fmt"
	"strings"

	"github.com/arch-couple/arch-couple-installer/pkg/arch_chroot"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Homepath string `json:"homepath"`
	Sudoer   bool   `json:"sudoer"`
}

func New(username, password, homepath string, sudoer bool) (*User, error) {
	if strings.TrimSpace(username) == "" || strings.TrimSpace(password) == "" {
		return nil, UserInstantiationError{
			err: errors.New("Can't create user with empty username or password"),
		}
	}

	if strings.TrimSpace(homepath) == "" {
		homepath = fmt.Sprintf("/home/%s", username)
	} else if homepath[0] != '/' {
		return nil, UserInstantiationError{
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

func CreateUsers(users []User) error {
	err := setupSudoerFile()
	if err != nil {
		return err
	}

	for _, user := range users {
		err = userAdd(user.Username, user.Homepath)
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
	}

	return nil
}

func addToSudoer(username string) error {
	addToWheel := fmt.Sprintf("usermod -aG wheel %s", username)

	err := arch_chroot.Run(addToWheel)
	if err != nil {
		return err
	}

	return nil
}

func userAdd(username, homepath string) error {
	createCommand := fmt.Sprintf("useradd -m %s -d %s", username, homepath)

	err := arch_chroot.Run(createCommand)
	if err != nil {
		return err
	}

	return nil
}

func setPassword(username, password string) error {
	command := fmt.Sprintf("echo %s | passwd %s -s", password, username)

	err := arch_chroot.Run(command)
	if err != nil {
		return err
	}

	return nil
}

func setupSudoerFile() error {
	wheelLine := "%wheel      ALL=(ALL:ALL) ALL"
	command := fmt.Sprintf("echo \"%s\" >> /etc/sudoers", wheelLine)

	err := arch_chroot.Run(command)
	if err != nil {
		return err
	}

	return nil
}
