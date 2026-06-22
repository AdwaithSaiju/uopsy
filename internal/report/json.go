package report

import (
	"encoding/json"
	"os"

	"github.com/AdwaithSaiju/uopsy/internal/models"
)

func WriteJSON(devices []models.USBDevice, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(devices)
}
