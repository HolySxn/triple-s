package utils

import (
	"encoding/csv"
	"io"
	"os"
)

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

func UpdateCSV(dir string, name string) {
}

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

func FindName(dir, name string) (bool, int) {
	file, err := os.Open(dir)
	if err != nil {
		return false, -1
	}
	defer file.Close()

	reader := csv.NewReader(file)

	_, err = reader.Read()
	if err != nil {
		return false, -1
	}

	index := 1
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}

		if record[0] == name {
			return true, index
		}

		index++
	}

	return false, -1
}

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
