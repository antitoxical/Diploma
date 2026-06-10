package deliverykind

import "strings"

const (
	Courier = "courier"
	Pickup  = "pickup"
	PVZ     = "pvz"
	Postal  = "postal"
)

func Normalize(kind string) string {
	switch strings.ToLower(strings.TrimSpace(kind)) {
	case Pickup, PVZ, Postal, Courier:
		return strings.ToLower(strings.TrimSpace(kind))
	default:
		return Courier
	}
}

func IsValid(kind string) bool {
	switch Normalize(kind) {
	case Courier, Pickup, PVZ, Postal:
		return true
	default:
		return false
	}
}
