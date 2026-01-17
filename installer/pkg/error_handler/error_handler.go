package error_handler

import (
	"errors"

	"github.com/october-os/october-installer/pkg/arch_chroot"
	"github.com/october-os/october-installer/pkg/core"
	"github.com/october-os/october-installer/pkg/grub"
	"github.com/october-os/october-installer/pkg/hostname"
	"github.com/october-os/october-installer/pkg/locale"
	"github.com/october-os/october-installer/pkg/mirrors"
	"github.com/october-os/october-installer/pkg/partition"
	"github.com/october-os/october-installer/pkg/postinstall"
	"github.com/october-os/october-installer/pkg/timezone"
	"github.com/october-os/october-installer/pkg/user"
)

// Returns the exit code that should be used depending on the error.
func GetExitCode(e error) int {
	var (
		archChrootErr  arch_chroot.ArchChrootError
		coreInstallErr core.CoreInstallError
		grubErr        grub.GrubError
		hostnameErr    hostname.HostnameError
		localeErr      locale.LocaleError
		mirrorListErr  mirrors.MirrorListError
		partitionErr   partition.PartitionError
		postInstallErr postinstall.PostInstallError
		timezoneErr    timezone.TimezoneError
		newUserErr     user.NewUserError
	)

	switch {
	case errors.As(e, &archChrootErr):
		return 1
	case errors.As(e, &coreInstallErr):
		return 2
	case errors.As(e, &grubErr):
		return 3
	case errors.As(e, &hostnameErr):
		return 4
	case errors.As(e, &localeErr):
		return 5
	case errors.As(e, &mirrorListErr):
		return 6
	case errors.As(e, &partitionErr):
		return 7
	case errors.As(e, &postInstallErr):
		return 8
	case errors.As(e, &timezoneErr):
		return 9
	case errors.As(e, &newUserErr):
		return 10
	default:
		return 67
	}
}
