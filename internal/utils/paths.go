package utils

import (
	"os"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func PahtFiles(path string) []string{
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()

	names, _ := f.Readdirnames(0)

	return names
}
