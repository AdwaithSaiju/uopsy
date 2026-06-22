//go:build darwin

package collector

import (
	"os/exec"
	"strings"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

func collect() ([]models.USBDevice, error) {
	out, err := exec.Command("system_profiler", "SPUSBDataType", "-xml").Output()
	if err != nil {
		return nil, err
	}

	return parseSystemProfiler(string(out)), nil
}

func parseSystemProfiler(xml string) []models.USBDevice {
	var devices []models.USBDevice
	lines := strings.Split(xml, "\n")
	var current models.USBDevice

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.Contains(line, "<key>vendor_id</key>") {
			current.VendorID = extractValue(line)
		} else if strings.Contains(line, "<key>product_id</key>") {
			current.ProductID = extractValue(line)
		} else if strings.Contains(line, "<key>_name</key>") {
			current.Product = extractValue(line)
		} else if strings.Contains(line, "<key>manufacturer</key>") {
			current.Manufacturer = extractValue(line)
		} else if strings.Contains(line, "<key>serial_num</key>") {
			current.Serial = extractValue(line)
		} else if strings.Contains(line, "</dict>") && current.VendorID != "" {
			devices = append(devices, current)
			current = models.USBDevice{}
		}
	}

	return devices
}

func extractValue(line string) string {
	start := strings.Index(line, "<string>")
	if start == -1 {
		return ""
	}
	start += len("<string>")
	end := strings.Index(line[start:], "</string>")
	if end == -1 {
		return ""
	}
	return line[start : start+end]
}
