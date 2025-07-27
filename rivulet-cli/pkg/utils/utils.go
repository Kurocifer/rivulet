package utils

import (
	"os"
)

// GetWorkDir, returns the working directory of the rivulet daemon
func GetWorkDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./rivulet/.daemon"
	}

	return homeDir + "/rivulet/.daemon"
}

// CreateWorkDir, creates teh rivulet working directory
func CreateWorkDir() error {
	dir := GetWorkDir()

	return os.MkdirAll(dir, 0755)
}
