package models

import "fmt"

type USBDevice struct {
	VendorID     string `json:"vendor_id"`
	ProductID    string `json:"product_id"`
	Product      string `json:"product"`
	Manufacturer string `json:"manufacturer"`
	Serial       string `json:"serial"`
	Bus          string `json:"bus"`
	Port         string `json:"port"`
}

func (d USBDevice) String() string {
	return fmt.Sprintf("%s:%s  %s  %s", d.VendorID, d.ProductID, d.Manufacturer, d.Product)
}
