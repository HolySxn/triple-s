package object

import (
	"net/http"
	"os"
)

func PutObject(data []byte, dir string) int {
	file, err := os.Create(dir)
	if err != nil{
		return http.StatusInternalServerError
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}
