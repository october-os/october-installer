package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/october-os/october-installer/pkg/core"
	"github.com/october-os/october-installer/pkg/error_handler"
	"github.com/october-os/october-installer/pkg/fstab"
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

	if *json == "" || json == nil {
		fmt.Println("Missing 'json' arg.")
		os.Exit(1)
	}

	flag.Parse()
	var installationConfig *json_parser.Installation

	if *json != "" {
		fmt.Println("Parsing JSON configuration...")

		var err error
		installationConfig, err = json_parser.ParseJson(*json)
		if err != nil {
			exitWithErrorCode(err, err.Error())
		}

		fmt.Println("JSON configuration parsed.")
	}

	preInstallStep(&installationConfig.Drives, &installationConfig.MirrorCountries)

	fmt.Println("Installing core Arch Linux...")
	if err := core.InstallBasicInstallation(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Core Arch Linux installation finished.")
	fmt.Println("Generating fstab...")
	if err := fstab.GenerateFstab(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Fstab generated.")
	fmt.Println("Setting up timezone...")

	setTime(installationConfig.Timezone)

	fmt.Println("Finished setting up timezone.")
	fmt.Println("Setting up locales...")

	if err := locale.GenerateLocales(installationConfig.Locale); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished setting up locales.")
	fmt.Println("Setting up hostname...")

	if err := hostname.SetHostname(installationConfig.Hostname); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished setting up hostname.")

	postInstallStep(installationConfig.BestEffortGpu)

	userCreation(&installationConfig.Users, installationConfig.RootPassword)

	fmt.Println("October Linux installation done!")
}

// Creates all the given users and sets the
// root password.
func userCreation(users *[]user.User, rootPassword string) {
	for _, userToCreate := range *users {
		fmt.Printf("Creating user %s...\n", userToCreate.Username)
		if err := user.CreateUser(&userToCreate); err != nil {
			exitWithErrorCode(err, err.Error())
		}

		fmt.Println("User created.")
	}

	fmt.Println("Setting up root password...")
	if err := user.SetRootPassword(rootPassword); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished setting up root password.")
}

// Runs all the pre install steps like doing partitions
// and setting up mirrors.
func preInstallStep(drives *[]partition.Drive, mirrorCountries *[]string) {
	fmt.Println("Setting up partitions...")

	if err := partition.SetupPartitions(*drives); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Partitions created.")
	fmt.Println("Setting up mirror list...")

	if err := mirrors.SetMirrorList(*mirrorCountries); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished setting up mirror list.")
}

// Runs all the post install steps like installing
// packages, GPU drivers etc.
func postInstallStep(installGpuDrivers bool) {
	fmt.Println("Setting up branding...")
	if err := postinstall.SetupBranding(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Branding set up.")
	fmt.Println("Enabling multilib repository...")

	if err := postinstall.EnableMultilibRepo(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Multilib repository enabled.")
	fmt.Println("Installing packages from official repositories...")

	if err := postinstall.InstallOfficialPackages(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished installing packages.")
	fmt.Println("Installing AUR helper and AUR packages...")

	if err := postinstall.InstallAurHelperAndPackages(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished installing AUR helper and packages.")

	if installGpuDrivers {
		fmt.Println("Checking and installing GPU drivers...")
		if err := postinstall.BestEffortGPUDrivers(); err != nil {
			exitWithErrorCode(err, err.Error())
		}

		fmt.Println("GPU drivers installed.")
	}

	fmt.Println("Setting up sudo...")

	if err := postinstall.SetupSudo(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished setting up sudo.")
	fmt.Println("Installing grub as bootloader...")

	if err := grub.InstallGrub(); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	fmt.Println("Finished setting grub as bootloader.")
}

// Runs the function to set the timezone
// and sync the hardware clock
func setTime(tmz string) {
	if err := timezone.SetTime(tmz); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := timezone.SetHwClock(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
}

// Exits the installer with the given exit code
// and prints the given message in STDERR.
func exitWithErrorCode(e error, m string) {
	exit_code := error_handler.GetExitCode(e)
	fmt.Fprintln(os.Stderr, m)
	os.Exit(exit_code)
}
