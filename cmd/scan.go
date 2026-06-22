package cmd

import (
	"fmt"

	"github.com/AdwaithSaiju/uopsy/internal/collector"
)

func Scan() {
	fmt.Println(" [+] Scanning USB devices...")

	devices, err := collector.GetCurrentOS()
	if err != nil {
		fmt.Printf(" [-] Error scanning devices: %v\n", err)
		return
	}

	if len(devices) == 0 {
		fmt.Println(" [!] No USB devices detected.")
		return
	}

	fmt.Printf(" [+] Found %d devices:\n\n", len(devices))
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
