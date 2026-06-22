package collector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

func GetLinuxDetails() ([]models.USBDevice, error) {
	basePath := "/sys/bus/usb/devices/"

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read usb directory: %w", err)
	}

	var unifiedDevices []models.USBDevice

	for _, entry := range entries {
		name := entry.Name()

		if !strings.Contains(name, "-") && !strings.HasPrefix(name, "usb") {
			continue
		}

		devPath := filepath.Join(basePath, name)

		vid := readSysfsFile(filepath.Join(devPath, "idVendor"))
		pid := readSysfsFile(filepath.Join(devPath, "idProduct"))
		serial := readSysfsFile(filepath.Join(devPath, "serial"))
		productName := readSysfsFile(filepath.Join(devPath, "product"))

		if vid == "" && pid == "" {
			continue
		}

		if productName == "" {
			productName = "Unknown USB Device"
		}
		if serial == "" {
			serial = "Unknown"
		}

		unifiedDevices = append(unifiedDevices, models.USBDevice{
			Name:      productName,
			VendorID:  vid,
			ProductID: pid,
			Serial:    serial,
			Connected: true, // sysfs paths only contain currently active devices
		})
	}

	return unifiedDevices, nil
}

func readSysfsFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
