package collector

import "runtime"

func GetCurrentOS() string {
	switch runtime.GOOS {
	case "windows":
		return "Windows"
	case "linux":
		return "Linux"
	case "darwin":
		return "Mac"
	default:
		return "Unknown"
	}
}
