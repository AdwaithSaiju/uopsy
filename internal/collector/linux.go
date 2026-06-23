package collector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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
			Connected: true,
		})
	}

	return unifiedDevices, nil
}

var (
	vendorRe   = regexp.MustCompile(`idVendor=([0-9a-fA-F]{4})`)
	productRe  = regexp.MustCompile(`idProduct=([0-9a-fA-F]{4})`)
	productNRe = regexp.MustCompile(`Product:\s*(.+)`)
	serialRe   = regexp.MustCompile(`SerialNumber:\s*(.+)`)
	devAddrRe  = regexp.MustCompile(`usb\s+([\d-]+\.?\d*)`)
)

func GetLinuxHistory() ([]models.USBDevice, error) {
	cmd := exec.Command("journalctl", "-k", "--no-pager")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed reading kernel logs: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	type logDev struct {
		addr        string
		vendorID    string
		productID   string
		productName string
		serial      string
		disconnected bool
		seen        bool
		seenAddr    bool
	}

	devMap := make(map[string]*logDev)
	var order []string

	for _, line := range lines {
		addrMatch := devAddrRe.FindStringSubmatch(line)
		if addrMatch == nil {
			continue
		}
		addr := addrMatch[1]

		if _, ok := devMap[addr]; !ok {
			devMap[addr] = &logDev{addr: addr}
			order = append(order, addr)
		}
		d := devMap[addr]

		if strings.Contains(line, "USB disconnect") {
			d.disconnected = true
			continue
		}

		if !d.seenAddr && strings.Contains(line, "New USB device found") {
			d.seenAddr = true
		}

		if v := vendorRe.FindStringSubmatch(line); v != nil {
			d.vendorID = v[1]
		}
		if p := productRe.FindStringSubmatch(line); p != nil {
			d.productID = p[1]
		}
		if p := productNRe.FindStringSubmatch(line); p != nil {
			d.productName = strings.TrimSpace(p[1])
		}
		if s := serialRe.FindStringSubmatch(line); s != nil {
			d.serial = strings.TrimSpace(s[1])
		}

		if strings.Contains(line, "New USB device found") {
			d.seen = true
		}
	}

	var devices []models.USBDevice
	for _, addr := range order {
		d := devMap[addr]
		if !d.seen {
			continue
		}
		if d.vendorID == "" && d.productID == "" {
			continue
		}

		name := d.productName
		if name == "" {
			name = "Unknown USB Device"
		}
		serial := d.serial
		if serial == "" {
			serial = "Unknown"
		}

		devices = append(devices, models.USBDevice{
			Name:      name,
			VendorID:  d.vendorID,
			ProductID: d.productID,
			Serial:    serial,
			Connected: !d.disconnected,
		})
	}

	return devices, nil
}

func readSysfsFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
