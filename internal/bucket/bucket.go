package bucket

import (
	"errors"
	"os"
	"time"

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

	metaData := []string{name, time.Now().Format("2006-01-02 15:04:05 MST"), time.Now().Format("2006-01-02 15:04:05 MST"), "Active"}
	err = utils.WriteCSV(dir+"/buckets.csv", metaData)
	if err != nil{
		return err
	}

	return nil
}
