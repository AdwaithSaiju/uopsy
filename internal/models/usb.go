package models

type USBDevice struct {
	Name      string
	VendorID  string
	ProductID string
	Serial    string
	Connected bool
}
