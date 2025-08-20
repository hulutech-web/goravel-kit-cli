package utils

import (
	"os"
)

func DirectoryExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func MoveDirectory(source, destination string) error {
	return os.Rename(source, destination)
}

func RemoveDirectory(path string) error {
	return os.RemoveAll(path)
}
