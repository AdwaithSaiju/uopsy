//go:build linux

package collector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

func collect() ([]models.USBDevice, error) {
	devices := []models.USBDevice{}

	entries, err := os.ReadDir("/sys/bus/usb/devices")
	if err != nil {
		return nil, err
	}

	for _, e := range entries {
		name := e.Name()
		if !strings.Contains(name, ":") {
			continue
		}

		dev := models.USBDevice{Bus: "usb", Port: name}

		if vid := readFile(filepath.Join("/sys/bus/usb/devices", name, "idVendor")); vid != "" {
			dev.VendorID = vid
		}
		if pid := readFile(filepath.Join("/sys/bus/usb/devices", name, "idProduct")); pid != "" {
			dev.ProductID = pid
		}
		if p := readFile(filepath.Join("/sys/bus/usb/devices", name, "product")); p != "" {
			dev.Product = p
		}
		if m := readFile(filepath.Join("/sys/bus/usb/devices", name, "manufacturer")); m != "" {
			dev.Manufacturer = m
		}
		if s := readFile(filepath.Join("/sys/bus/usb/devices", name, "serial")); s != "" {
			dev.Serial = s
		}

		devices = append(devices, dev)
	}

	return devices, nil
}

func readFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
