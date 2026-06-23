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

type winRegDevice struct {
	VID         string
	PID         string
	Serial      string
	Description string
}

func GetWindowsDetails() ([]models.USBDevice, error) {
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

func GetWindowsHistory() ([]models.USBDevice, error) {
	psCommand := `
$regPath = 'HKLM:\SYSTEM\CurrentControlSet\Enum\USB'
$devices = @()
Get-ChildItem $regPath -ErrorAction SilentlyContinue | ForEach-Object {
	$vidPid = $_.PSChildName
	$parts = $vidPid -split '&'
	$vid = ''
	$pid = ''
	foreach ($p in $parts) {
		if ($p -match '^VID_(.+)') { $vid = $matches[1] }
		if ($p -match '^PID_(.+)') { $pid = $matches[1] }
	}
	Get-ChildItem $_.PSPath -ErrorAction SilentlyContinue | ForEach-Object {
		$serial = $_.PSChildName
		$desc = ($_.GetValue('DeviceDesc') -replace '^.*;', '')
		if (-not $desc) { $desc = 'Unknown USB Device' }
		if ($serial -eq '') { $serial = 'Unknown' }
		$devices += @{VID=$vid; PID=$pid; Serial=$serial; Description=$desc}
	}
}
$devices | ConvertTo-Json -Compress
`
	cmd := exec.Command("powershell", "-Command", psCommand)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed reading USB registry: %w", err)
	}

	if len(output) == 0 || strings.TrimSpace(string(output)) == "null" {
		return []models.USBDevice{}, nil
	}

	var regDevices []winRegDevice
	if err := json.Unmarshal(output, &regDevices); err != nil {
		return nil, fmt.Errorf("failed parsing registry JSON: %w", err)
	}

	devices := make([]models.USBDevice, 0, len(regDevices))
	for _, rd := range regDevices {
		if rd.VID == "" && rd.PID == "" {
			continue
		}
		devices = append(devices, models.USBDevice{
			Name:      rd.Description,
			VendorID:  rd.VID,
			ProductID: rd.PID,
			Serial:    rd.Serial,
			Connected: false,
		})
	}

	return devices, nil
}

func parseInstanceID(instanceID string) (vid, pid, serial string) {
	if instanceID == "" {
		return "Unknown", "Unknown", "Unknown"
	}

	parts := strings.Split(instanceID, "\\")
	if len(parts) < 3 {
		return "Unknown", "Unknown", instanceID
	}

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
