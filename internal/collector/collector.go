package collector

import (
	"fmt"
	"os"
	"runtime"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

func GetCurrentOS() ([]models.USBDevice, error) {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsDetails()
	case "linux":
		return GetLinuxDetails()
	default:
		fmt.Printf(" [-] Operating system %s is unsupported.\n", runtime.GOOS)
		os.Exit(1)
		return nil, nil
	}
}

func GetUSBHistory() ([]models.USBDevice, error) {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsHistory()
	case "linux":
		return GetLinuxHistory()
	default:
		return nil, nil
	}
}

func GetAllDevices() ([]models.USBDevice, error) {
	current, err := GetCurrentOS()
	if err != nil {
		return nil, err
	}

	history, err := GetUSBHistory()
	if err != nil {
		return nil, err
	}

	seen := make(map[string]bool)
	for _, d := range current {
		key := d.VendorID + ":" + d.ProductID + ":" + d.Serial
		seen[key] = true
	}

	for _, d := range history {
		key := d.VendorID + ":" + d.ProductID + ":" + d.Serial
		if !seen[key] {
			current = append(current, d)
			seen[key] = true
		}
	}

	return current, nil
}
