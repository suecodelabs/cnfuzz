package util

// IsValidPort check if a port is valid
func IsValidPort(port int) bool {
	if port <= 0 || port > 65535 {
		return false
	}
	return true
}
