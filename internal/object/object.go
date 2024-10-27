package object

import (
	"io"
	"net/http"
	"os"
	"path"

	"triple-s/internal/utils"
)

func PutObject(data []byte, dir string) (int, string) {
	file, err := os.Create(dir)
	if err != nil {
		return http.StatusInternalServerError, "Internal Server Error"
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return http.StatusInternalServerError, "internal Server Error"
	}

	return http.StatusOK, "Object was successfully added"
}

func GetObject(dir string) ([]byte, int) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	data, err := io.ReadAll(file)
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	return data, http.StatusOK
}

func DeleteObject(name, dir string) (int, string) {
	if flag, index, _ := utils.FindName(path.Join(dir, "objects.csv"), name); flag {
		err := os.Remove(path.Join(dir, name))
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}

		err = utils.UpdateCSV(path.Join(dir, "objects.csv"), "delete", index, nil)
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}

		return http.StatusOK, "Object was successfully deleted"
	} else {
		return http.StatusNotFound, "Object is not found"
	}
}
