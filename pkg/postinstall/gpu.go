package postinstall

import (
	"errors"
	"strings"

	"github.com/jaypipes/ghw"
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

// Fetches the system's GPU information using ghw and returns it
func getGPUInfo() (gpuInfo, error) {
	info, err := ghw.GPU()
	if err != nil {
		return gpuInfo{}, err
	}

	if len(info.GraphicsCards) == 0 {
		return gpuInfo{}, errors.New("no graphics cards found")
	}

	card := info.GraphicsCards[0]
	if strings.Contains(card.DeviceInfo.Vendor.Name, "Advanced Micro Devices") {
		return gpuInfo{brand: "AMD", family: ""}, nil
	} else if strings.Contains(card.DeviceInfo.Vendor.Name, "Intel") {
		return gpuInfo{brand: "Intel", family: ""}, nil
	} else if strings.Contains(card.DeviceInfo.Vendor.Name, "NVIDIA") {
		productFamily := strings.Split(card.DeviceInfo.Product.Name, " ")
		return gpuInfo{brand: "NVIDIA", family: productFamily[0]}, nil
	}

	return gpuInfo{}, errors.New("no supported GPU device found")
}
