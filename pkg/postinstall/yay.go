package postinstall

import (
	"fmt"

	"github.com/october-os/october-installer/pkg/arch_chroot"
)

const yayPath string = "/tmp/yay-bin"

// Clones the repo of yay-bin, from the AUR, inside /tmp
// and installs it with the builder account.
func installYay() error {
	cloneYayBinary := fmt.Sprintf("git clone https://aur.archlinux.org/yay-bin.git %s", yayPath)
	chownToNobody := fmt.Sprintf("chown %s %s", username, yayPath)
	goInDir := fmt.Sprintf("cd %s", yayPath)
	makePkg := fmt.Sprintf("sudo -u %s makepkg -si --noconfirm", username)

	command := fmt.Sprintf("%s && %s && %s && %s", cloneYayBinary, chownToNobody, goInDir, makePkg)
	return arch_chroot.Run(command)
}
