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
	if strings.TrimSpace(username) == "" {
		return nil, errors.New("Can't create user with empty username")
	}

	if strings.TrimSpace(homepath) == "" {
		homepath = fmt.Sprintf("/home/%s", username)
	}

	return &User{
		Username: username,
		Password: password,
		Homepath: homepath,
		Sudoer:   sudoer,
	}, nil
}

func set_password(username, password string) error {
	command := fmt.Sprintf("echo %s | passwd %s -s", password, username)

	err := arch_chroot.Run(command)
	if err != nil {
		return err
	}

	return nil
}
