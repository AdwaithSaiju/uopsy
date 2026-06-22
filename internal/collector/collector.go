package collector

import (
	"fmt"
	"os"
	"runtime"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

func GetCurrentOS() ([]models.USBDevice, error) {
	switch runtime.GOOS {
	case "windows":
		return GetWindowsDetails()
	case "linux":
		return GetLinuxDetails()
	case "darwin":
		return GetMacOSDetails()
	default:
		fmt.Printf(" [-] Operating system %s is unsupported.\n", runtime.GOOS)
		os.Exit(1)
		return nil, nil
	}
}
