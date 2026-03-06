package postinstall

import (
	"fmt"
	"io"
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
		return &PostInstallError{err: err}
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
		if err := addPackages(officialPackagesFilePath, officialPackages); err != nil {
			return &PostInstallError{err: err}
		}
	}
	if len(aurPackages) > 0 {
		if err := addPackages(aurPackagesFilePath, aurPackages); err != nil {
			return &PostInstallError{err: err}
		}
	}
	return nil
}

// Fetches the system's GPU information using lspci and returns it
func getGPUInfo() (gpuInfo, error) {
	command := "lspci | grep -i 'VGA compatible controller'"
	cmd := exec.Command("/bin/bash", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return gpuInfo{}, err
	}

	if err := cmd.Start(); err != nil {
		return gpuInfo{}, err
	}
	var stdoutOutput []byte
	if stdoutOutput, err = io.ReadAll(stdout); err != nil {
		return gpuInfo{}, err
	}
	if err := cmd.Wait(); err != nil {
		return gpuInfo{}, err
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

	return gpuInfo{}, fmt.Errorf("error getting GPU brand: not found")
}
