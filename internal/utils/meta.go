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

func WriteCSV(dir string, record []string) error {
	file, err := os.OpenFile(dir, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write(record)
	if err != nil {
		return err
	}
	return nil
}

func UpdateCSV(dir string, name string) {
}
