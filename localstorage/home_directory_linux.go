package localstorage

import "os"

func homeDirectory() string {
	// On Linux, the home directory is typically found in the $HOME environment variable.
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/home/default"
	}
	return homeDir
}
