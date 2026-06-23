package cmd

import (
	"fmt"

	"github.com/AdwaithSaiju/uopsy/internal/collector"
)

func Scan() {
	fmt.Println(" [+] Scanning USB devices...")

	devices, err := collector.GetAllDevices()
	if err != nil {
		fmt.Printf(" [-] Error scanning devices: %v\n", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println(" [!] No USB devices detected.")
		return
	}

	connected := 0
	historical := 0
	for _, dev := range devices {
		if dev.Connected {
			connected++
		} else {
			historical++
		}
	}

	fmt.Printf(" [+] Found %d devices (%d connected, %d historical):\n\n", len(devices), connected, historical)
	for _, dev := range devices {
		status := "DISCONNECTED"
		if dev.Connected {
			status = "CONNECTED"
		}

		fmt.Printf(" [%s] %s\n", status, dev.Name)
		fmt.Printf("   ├─ VID:    %s\n", dev.VendorID)
		fmt.Printf("   ├─ PID:    %s\n", dev.ProductID)
		fmt.Printf("   └─ Serial: %s\n\n", dev.Serial)
	}
}
