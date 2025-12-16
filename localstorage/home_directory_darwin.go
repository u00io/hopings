package localstorage

import "os"

func homeDirectory() string {
	// MacOS typically uses the $HOME environment variable to determine the user's home directory.
	homeDir := os.Getenv("HOME")
	if homeDir == "" {
		homeDir = "/"
	}
	return homeDir
}
