package object

import (
	"io"
	"net/http"
	"os"
	"path"
	"triple-s/internal/utils"
)

// PutObject stores the given data as an object at the specified directory path.
// It creates a new file, writes the data to it, and returns the HTTP status code and a message.
func PutObject(data []byte, dir string) (int, string) {
	// Create the file at the specified path
	file, err := os.Create(dir)
	if err != nil {
		return http.StatusInternalServerError, "Internal Server Error"
	}
	defer file.Close()

	// Write the provided data to the file
	_, err = file.Write(data)
	if err != nil {
		return http.StatusInternalServerError, "internal Server Error"
	}

	// Return success status and message if writing was successful
	return http.StatusOK, "Object was successfully added"
}

// GetObject retrieves the data of an object from the specified directory path.
// It opens the file, reads its content, and returns the data along with the HTTP status code.
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

// DeleteObject deletes an object with the specified name from the directory.
// It checks for the object in the metadata file ("objects.csv") before deleting it from disk
// and removing its metadata record.
func DeleteObject(name, dir string) (int, string) {
	// Check if the object exists by searching for it in "objects.csv"
	if flag, index, _ := utils.FindName(path.Join(dir, "objects.csv"), name); flag {
		// Remove the object file from the directory
		err := os.Remove(path.Join(dir, name))
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}

		// Remove the object's entry from "objects.csv"
		err = utils.UpdateCSV(path.Join(dir, "objects.csv"), "delete", index, nil)
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}

		return http.StatusOK, "Object was successfully deleted"
	} else {
		return http.StatusNotFound, "Object is not found"
	}
}
