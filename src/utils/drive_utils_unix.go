//go:build !windows
// +build !windows

package utils

func GetLocalDrives() ([]string, error) {
	// On Unix-like systems, return the root directory
	return []string{"/"}, nil
}
