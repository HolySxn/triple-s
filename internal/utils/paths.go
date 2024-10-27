package utils

import (
	"os"
)
// IsExist check if the provided path is exist or not
func IsExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
