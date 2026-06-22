package cmd

import (
	"fmt"

	"github.com/AdwaithSaiju/uopsy/internal/collector"
)

func Scan() {
	currentOS := collector.GetCurrentOS()

	fmt.Printf(" [+] Detected OS: %s\n", currentOS)
	fmt.Println(" [+] Scanning USB devices...")

}
