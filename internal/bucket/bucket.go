package bucket

import (
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

	return nil
}
