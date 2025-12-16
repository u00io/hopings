package localstorage

import "os"

func homeDirectory() string {
	homeDir := os.Getenv("USERPROFILE")
	if homeDir == "" {
		homeDir = "C:\\Users\\Default"
	}
	return homeDir
}
