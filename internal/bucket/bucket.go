package bucket

import (
	"encoding/csv"
	"errors"
	"os"

	"triple-s/internal/utils"
)

func CreateBucket(name string, dir string) error {
	if !utils.IsValidBucketName(name) {
		return errors.New("invalid bucket name")
	}

	bucket_dir := dir + "/" + name
	err := os.Mkdir(bucket_dir, os.ModePerm)
	if err != nil {
		return err
	}

	createCSVbucket(bucket_dir)
	return nil
}

func createCSVbucket(name string) error {
	file, err := os.Create(name + "/objects.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"Name", "CreationTime", "LastModifiedTime", "Status"}
	if err := writer.Write(header); err != nil {
		return err
	}

	return nil
}
