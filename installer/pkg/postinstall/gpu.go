package postinstall

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"slices"
	"strings"
)

// Packages for each brand
// AMD: mesa, lib32-mesa, vulkan-radeon, lib32-vulkan-radeon
// Intel: mesa, lib32-mesa, vulkan-intel, lib32-vulkan-intel
// NVIDIA: TUXXX, GAXXX, ADXXX: nvidia-open ---- GMXXX, GPXXX, GVXXX: nvidia-580xx-dkms (aur)

var amdGPUPackages []string = []string{"mesa", "lib32-mesa", "vulkan-radeon", "lib32-vulkan-radeon"}
var intelGPUPackages []string = []string{"mesa", "lib32-mesa", "vulkan-intel", "lib32-vulkan-intel"}

const nvidiaOpenGPUPackage string = "nvidia-open"
const nvidiaProprietaryGPUPackage string = "nvidia-580xx-dkms"

// (newest) NVIDIA GPU families (https://nouveau.freedesktop.org/CodeNames.html)
var nvidiaGPUFamilies []string = []string{"TU", "GA", "AD", "GM", "GP", "GV"}

// GPUInfo represents a GPU's fetched information
type gpuInfo struct {
	brand  string
	family string
}

// Chooses which packages should be installed for the system's GPU information
// Can return error type: PostInstallError
func BestEffortGPUDrivers() error {
	gpuInfo, err := getGPUInfo()
	if err != nil {
		return err
	}

	var officialPackages []string
	var aurPackages []string

	switch gpuInfo.brand {
	case "Intel":
		officialPackages = append(officialPackages, intelGPUPackages...)
	case "AMD":
		officialPackages = append(officialPackages, amdGPUPackages...)
	case "NVIDIA":
		switch gpuInfo.family[:2] {
		case "TU", "GA", "AD":
			officialPackages = append(officialPackages, nvidiaOpenGPUPackage)
		case "GM", "GP", "GV":
			aurPackages = append(aurPackages, nvidiaProprietaryGPUPackage)
		}
	}

	if len(officialPackages) > 0 {
		addGPUPackages(officialPackagesFilePath, officialPackages)
	}
	if len(aurPackages) > 0 {
		addGPUPackages(aurPackagesFilePath, aurPackages)
	}
	return nil
}

// Fetches the system's GPU information using lspci and returns it
// Can return error type: PostInstallError
func getGPUInfo() (gpuInfo, error) {
	command := "lspci | grep -i 'VGA compatible controller'"
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return gpuInfo{}, &PostInstallError{
			err: fmt.Errorf("error piping stdout: error=%s", err.Error()),
		}
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return gpuInfo{}, &PostInstallError{
			err: fmt.Errorf("error piping stderr: error=%s", err.Error()),
		}
	}
	if err := cmd.Start(); err != nil {
		stderrOutput, _ := io.ReadAll(stderr)
		return gpuInfo{}, &PostInstallError{
			err: fmt.Errorf("error getting GPU information: error=%s", string(stderrOutput)),
		}
	}
	var stdoutOutput []byte
	if stdoutOutput, err = io.ReadAll(stdout); err != nil {
		return gpuInfo{}, &PostInstallError{
			err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
		}
	}
	if err := cmd.Wait(); err != nil {
		return gpuInfo{}, &PostInstallError{
			err: fmt.Errorf("error reading stdout: error=%s", err.Error()),
		}
	}

	stdoutOutputString := string(stdoutOutput)
	if strings.Contains(stdoutOutputString, "Intel") {
		return gpuInfo{
			brand: "Intel",
		}, nil
	}
	if strings.Contains(stdoutOutputString, "AMD") {
		return gpuInfo{
			brand: "AMD",
		}, nil
	}
	if strings.Contains(stdoutOutputString, "NVIDIA") {
		for p := range strings.SplitSeq(stdoutOutputString, " ") {
			if len(p) == 5 && slices.Contains(nvidiaGPUFamilies, p[:2]) {
				return gpuInfo{
					brand:  "NVIDIA",
					family: p,
				}, nil
			}
		}
	}

	return gpuInfo{}, &PostInstallError{
		err: fmt.Errorf("error getting GPU brand: not found"),
	}
}

// Adds the packages to a given file path in this format: "- package\n"
// Can return error type: PostInstallError
func addGPUPackages(filePath string, packages []string) error {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0)
	if err != nil {
		return &PostInstallError{
			err: fmt.Errorf("error writing gpu packages to '%s': error=%s", filePath, err.Error()),
		}
	}
	defer file.Close()
	for _, p := range packages {
		fmt.Fprintf(file, "- %s\n", p)
	}
	return nil
}
