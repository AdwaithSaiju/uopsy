package collector

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

type winPnpDevice struct {
	InstanceId   string `json:"InstanceId"`
	FriendlyName string `json:"FriendlyName"`
	Class        string `json:"Class"`
	Present      bool   `json:"Present"`
}

func GetWindowsDetails() ([]models.USBDevice, error) {
	// @(...) forces PowerShell to always return a JSON array even if only 1 device is found
	psCommand := "@(Get-PnpDevice -Class 'USB' | Select-Object InstanceId, FriendlyName, Class, Present) | ConvertTo-Json -Compress"

	cmd := exec.Command("powershell", "-Command", psCommand)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed executing powershell scanner: %w", err)
	}

	if len(output) == 0 || strings.TrimSpace(string(output)) == "null" {
		return []models.USBDevice{}, nil
	}

	var pnpDevices []winPnpDevice
	if err := json.Unmarshal(output, &pnpDevices); err != nil {
		return nil, fmt.Errorf("failed parsing PnpDevice JSON: %w", err)
	}

	unifiedDevices := make([]models.USBDevice, 0, len(pnpDevices))

	for _, pnp := range pnpDevices {
		name := pnp.FriendlyName
		if name == "" {
			name = "Unknown USB Device"
		}

		vid, pid, serial := parseInstanceID(pnp.InstanceId)

		unifiedDevices = append(unifiedDevices, models.USBDevice{
			Name:      name,
			VendorID:  vid,
			ProductID: pid,
			Serial:    serial,
			Connected: pnp.Present,
		})
	}

	return unifiedDevices, nil
}

func parseInstanceID(instanceID string) (vid, pid, serial string) {
	if instanceID == "" {
		return "Unknown", "Unknown", "Unknown"
	}

	parts := strings.Split(instanceID, "\\")
	if len(parts) < 3 {
		return "Unknown", "Unknown", instanceID
	}

	// Extracts VID/PID from layout: VID_XXXX&PID_XXXX
	ids := strings.Split(parts[1], "&")
	for _, id := range ids {
		if strings.HasPrefix(id, "VID_") {
			vid = strings.TrimPrefix(id, "VID_")
		} else if strings.HasPrefix(id, "PID_") {
			pid = strings.TrimPrefix(id, "PID_")
		}
	}

	serial = parts[2]

	if vid == "" {
		vid = "Unknown"
	}
	if pid == "" {
		pid = "Unknown"
	}
	if serial == "" {
		serial = "Unknown"
	}

	return vid, pid, serial
}
