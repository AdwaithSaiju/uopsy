package collector

import "github.com/AdwaithSaiju/uopsy/internal/models"

func Collect() ([]models.USBDevice, error) {
	return collect()
}
