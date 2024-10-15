package utils

import (
	"encoding/csv"
	"os"
)

func CreateCSV(dir string, name string, header []string) error {
	file, err := os.Create(dir + "/" + name + ".csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if err := writer.Write(header); err != nil {
		return err
	}

	return nil
}
