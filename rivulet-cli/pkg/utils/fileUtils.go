package utils

import (
	"os"
)

func GetWorkDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "./rivulet/.daemon"
	}

	return homeDir + "/rivulet/.daemon"
}

func CreateWorkDir() error {
	dir := GetWorkDir()

	return os.MkdirAll(dir, 0755)
}
