package utils

import (
	"encoding/csv"
	"errors"
	"io"
	"os"
)

// CreateCSV creates a new CSV file at the specified path.
// It initializes the file and prepares it for writing, but does not write any data.
func CreateCSV(dir string) error {
	file, err := os.Create(dir)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	return nil
}

// WriteCSV appends a single record (row) to an existing CSV file.
func WriteCSV(dir string, record []string) error {
	file, err := os.OpenFile(dir, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o0644)
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

// WriteAllCSV writes multiple records (rows) to a CSV file in one batch.
func WriteAllCSV(dir string, record [][]string) error {
	file, err := os.OpenFile(dir, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.WriteAll(record)
	if err != nil {
		return err
	}
	return nil
}

// ReadCSV reads all data from a CSV file and returns it as a 2D slice of strings.
func ReadCSV(dir string) ([][]string, error) {
	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

// UpdateCSV updates or deletes a record in the CSV file based on the flag.
// The "delete" flag removes the record at the specified index, and "update" modifies it.
func UpdateCSV(dir string, flag string, index int, record []string) error {
	data, err := ReadCSV(dir)
	if err != nil {
		return err
	}

	err = CreateCSV(dir)
	if err != nil {
		return err
	}

	switch flag {
	case "delete":
		data = append(data[0:index], data[index+1:]...)
	case "update":
		data[index] = record
	default:
		return errors.New("not appropriate flag")
	}

	err = WriteAllCSV(dir, data)
	if err != nil {
		return err
	}

	return nil
}

// CreateStorage initializes storage by creating a directory and "buckets.csv" metadata file.
func CreateStorage(dir string) error {
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}

	// Create buckets.csv if not exist
	if !IsExist(dir + "/buckets.csv") {
		err = CreateCSV(dir + "/buckets.csv")
		if err != nil {
			return err
		}
		err = WriteCSV(dir+"/buckets.csv", []string{"Name", "CreationTime", "LastModifiedTime", "Status"})
		if err != nil {
			return err
		}
	}

	return nil
}

// FindName searches for a specific record in a CSV file by name (first column).
// It returns a flag indicating if the record is found, its index, and the record itself.
func FindName(dir, name string) (bool, int, []string) {
	file, err := os.Open(dir)
	if err != nil {
		return false, -1, nil
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return false, -1, nil
	}

	index := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if record[0] == name {
			return true, index, record
		}

		index++
	}

	return false, -1, nil
}

// IsEmptyCSV checks if a CSV file contains only the header row.
func IsEmptyCSV(dir string) bool {
	data, err := ReadCSV(dir)
	if err != nil {
		return false
	}

	if len(data) > 1 {
		return false
	}

	return true
}
