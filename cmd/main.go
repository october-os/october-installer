package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/october-os/october-installer/pkg/core"
	"github.com/october-os/october-installer/pkg/error_handler"
	"github.com/october-os/october-installer/pkg/fstab"
	"github.com/october-os/october-installer/pkg/grub"
	"github.com/october-os/october-installer/pkg/hostname"
	"github.com/october-os/october-installer/pkg/json_parser"
	"github.com/october-os/october-installer/pkg/keyring"
	"github.com/october-os/october-installer/pkg/locale"
	"github.com/october-os/october-installer/pkg/mirrors"
	"github.com/october-os/october-installer/pkg/partition"
	"github.com/october-os/october-installer/pkg/postinstall"
	"github.com/october-os/october-installer/pkg/timezone"
	"github.com/october-os/october-installer/pkg/user"
)

func main() {
	json := flag.String("json-raw", "", "Installation configuration as JSON")
	jsonFile := flag.String("json", "", "Path to installation configuration JSON file")

	flag.Parse()

	fmt.Println("Parsing JSON configuration...")
	installationConfig, err := getInstallConfiguration(json, jsonFile)
	if err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("JSON configuration parsed.")

	preInstallStep(&installationConfig.Drives, &installationConfig.MirrorCountries)

	fmt.Println("Installing core Arch Linux...")
	if err := core.InstallBasicInstallation(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Core Arch Linux installation finished.")

	fmt.Println("Mounting additional partitions...")
	if err := partition.MountAdditionalPartitions(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Additional partitions mounted.")

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

	userCreation(&installationConfig.Users, installationConfig.RootPassword)

	postInstallStep(installationConfig.BestEffortGpu, installationConfig.Users, installationConfig.ExtraPackages)

	fmt.Println("October Linux installation done!")
}

// getInstallConfiguration takes the inputs from the -json and -json-file flags and returns the
// installation configuration after being parsed or an error.
func getInstallConfiguration(json, jsonFile *string) (*json_parser.Installation, error) {
	if *json == "" && *jsonFile == "" {
		return nil, errors.New("missing 'json' or 'json-raw' arg")
	}

	var installationConfig *json_parser.Installation
	var err error

	if *json != "" {
		installationConfig, err = json_parser.ParseJson(*json)
		if err != nil {
			return nil, err
		}
	} else if *jsonFile != "" {
		content, err := os.ReadFile(*jsonFile)
		if err != nil {
			return nil, err
		}

		installationConfig, err = json_parser.ParseJson(string(content))
		if err != nil {
			return nil, err
		}
	}

	return installationConfig, nil
}

// userCreation creates all the given users and sets the
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

// preInstallStep runs all the pre-installation steps like doing partitions
// and setting up mirrors.
func preInstallStep(drives *[]partition.Drive, mirrorCountries *[]string) {
	fmt.Println("Updating keyring...")
	if err := keyring.UpdateKeyRing(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Keyring updated.")

	setupPartitionsPreInstall(drives)

	fmt.Println("Setting up mirror list...")
	if err := mirrors.SetMirrorList(*mirrorCountries); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Finished setting up mirror list.")
}

// setupPartitionsPreInstall creates and formats all partitions and mounts the
// system partitions needed for the core installation
func setupPartitionsPreInstall(drives *[]partition.Drive) {
	fmt.Println("Creating partitions...")
	if err := partition.CreatePartitions(*drives); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Partitions created.")

	fmt.Println("Formatting partitions...")
	if err := partition.FormatPartitions(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Partitions formatted.")

	fmt.Println("Mounting system partitions...")
	if err := partition.MountSystemPartitions(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("System partitions mounted.")
}

// postInstallStep runs all the post install steps like installing
// packages, GPU drivers etc.
func postInstallStep(installGpuDrivers bool, users []user.User, extraPackages postinstall.ExtraPackages) {
	fmt.Println("Enabling multilib repository...")
	if err := postinstall.EnableMultilibRepo(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Multilib repository enabled.")

	if installGpuDrivers {
		fmt.Println("Checking and installing GPU drivers...")
		if err := postinstall.BestEffortGPUDrivers(); err != nil {
			exitWithErrorCode(err, err.Error())
		}
		fmt.Println("GPU drivers installed.")
	}

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

	if len(extraPackages.OfficialRepositories) > 0 || len(extraPackages.AUR) > 0 {
		fmt.Println("Installing user-defined extra packages...")
		if err := postinstall.AddExtraPackages(extraPackages); err != nil {
			exitWithErrorCode(err, err.Error())
		}
		fmt.Println("Extra packages installed.")
	}

	fmt.Println("Setting up branding...")
	if err := postinstall.SetupBranding(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Branding set up.")

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

	fmt.Println("Installing October Linux configuration for all users...")
	if err := postinstall.AddConfigForUsers(users); err != nil {
		exitWithErrorCode(err, err.Error())
	}
	fmt.Println("Finished installing configuration for all users.")
}

// setTime runs the function to set the timezone
// and sync the hardware clock
func setTime(tmz string) {
	if err := timezone.SetTime(tmz); err != nil {
		exitWithErrorCode(err, err.Error())
	}

	if err := timezone.SetHwClock(); err != nil {
		exitWithErrorCode(err, err.Error())
	}
}

// exitWithErrorCode exits the installer with the given exit code
// and prints the given message in STDERR.
func exitWithErrorCode(e error, m string) {
	exitCode := error_handler.GetExitCode(e)
	fmt.Fprintln(os.Stderr, m)
	fmt.Println("it is recommended to reboot the machine before trying again")
	os.Exit(exitCode)
}
