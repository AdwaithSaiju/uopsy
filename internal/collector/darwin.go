package collector

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

type macUSBDevice struct {
	Name      string         `json:"_name"`
	VendorID  string         `json:"idVendor"`
	ProductID string         `json:"idProduct"`
	Serial    string         `json:"serial_num"`
	SubItems  []macUSBDevice `json:"_items"`
}

type macRoot struct {
	USBDataType []macUSBDevice `json:"SPUSBDataType"`
}

func GetMacOSDetails() ([]models.USBDevice, error) {
	cmd := exec.Command("system_profiler", "SPUSBDataType", "-json")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed executing system_profiler: %w", err)
	}

	var root macRoot
	if err := json.Unmarshal(output, &root); err != nil {
		return nil, fmt.Errorf("failed parsing macOS USB JSON: %w", err)
	}

	var unifiedDevices []models.USBDevice
	for _, topLevelItem := range root.USBDataType {
		traverseUSBTree(topLevelItem, &unifiedDevices)
	}

	return unifiedDevices, nil
}

func traverseUSBTree(item macUSBDevice, result *[]models.USBDevice) {
	if item.VendorID != "" || item.ProductID != "" {
		name := item.Name
		if name == "" {
			name = "Unknown USB Device"
		}

		serial := item.Serial
		if serial == "" {
			serial = "Unknown"
		}

		vid := strings.TrimSpace(strings.TrimPrefix(item.VendorID, "0x"))
		pid := strings.TrimSpace(strings.TrimPrefix(item.ProductID, "0x"))

		*result = append(*result, models.USBDevice{
			Name:      name,
			VendorID:  vid,
			ProductID: pid,
			Serial:    serial,
			Connected: true,
		})
	}

	for _, subItem := range item.SubItems {
		traverseUSBTree(subItem, result)
	}
}
