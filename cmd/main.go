package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/october-os/october-installer/pkg/core"
	"github.com/october-os/october-installer/pkg/error_handler"
	"github.com/october-os/october-installer/pkg/grub"
	"github.com/october-os/october-installer/pkg/hostname"
	"github.com/october-os/october-installer/pkg/json_parser"
	"github.com/october-os/october-installer/pkg/locale"
	"github.com/october-os/october-installer/pkg/mirrors"
	"github.com/october-os/october-installer/pkg/partition"
	"github.com/october-os/october-installer/pkg/postinstall"
	"github.com/october-os/october-installer/pkg/timezone"
	"github.com/october-os/october-installer/pkg/user"
)

func main() {
	json := flag.String("json", "", "Installation configuration JSON")

	flag.Parse()
	var installationConfig *json_parser.Installation

	if *json != "" {
		var err error
		installationConfig, err = json_parser.ParseJson(*json)
		if err != nil {
			exitWithErrorCode(err, err.Error())
		}
	}

	preInstallStep(&installationConfig.Drives, &installationConfig.MirrorCountries)

	if err := core.InstallBasicInstallation(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	setTime(installationConfig.Timezone)

	if err := locale.GenerateLocales(installationConfig.Locale); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := hostname.SetHostname(installationConfig.Hostname); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	postInstallStep(installationConfig.BestEffortGpu)

	userCreation(&installationConfig.Users, installationConfig.RootPassword)
}

func userCreation(users *[]user.User, rootPassword string) {
	for _, userToCreate := range *users {
		if err := user.CreateUser(&userToCreate); err != nil {
			exitWithErrorCode(err, err.Error())
		}
	}

	if err := user.SetRootPassword(rootPassword); err != nil {
		exitWithErrorCode(err, err.Error())
	}
}

func preInstallStep(drives *[]partition.Drive, mirrorCountries *[]string) {
	if err := partition.SetupPartitions(*drives); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := mirrors.SetMirrorList(*mirrorCountries); err != nil {
		exitWithErrorCode(err, err.Error())
	}
}

func postInstallStep(installGpuDrivers bool) {
	if err := postinstall.EnableMultilibRepo(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := postinstall.InstallOfficialPackages(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := postinstall.InstallAurHelperAndPackages(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if installGpuDrivers {
		if err := postinstall.BestEffortGPUDrivers(); err != nil {
			exitWithErrorCode(err, err.Error())
		}
	}

	if err := postinstall.SetupSudo(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := grub.InstallGrub(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
}

func setTime(tmz string) {
	if err := timezone.SetTime(tmz); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := timezone.SetHwClock(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
}

func exitWithErrorCode(e error, m string) {
	exit_code := error_handler.GetExitCode(e)
	fmt.Fprintln(os.Stderr, m)
	os.Exit(exit_code)
}
