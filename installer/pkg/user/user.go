package user

import (
	"fmt"

	"github.com/arch-couple/arch-couple-installer/pkg/arch_chroot"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Homepath string `json:"homepath"`
	Sudoer   bool   `json:"sudoer"`
}

func New(username, password, homepath string, sudoer bool) *User {
	// TODO: do all the checking and important stuff
}

func set_password(username, password string) error {
	command := fmt.Sprintf("echo %s | passwd %s -s", password, username)

	err := arch_chroot.Run(command)
	if err != nil {
		return err
	}

	return nil
}
