package bucket

import (
	"encoding/csv"
	"encoding/xml"
	"errors"
	"os"
	"time"

	"triple-s/internal/utils"
)

type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Bucket struct {
	Name             string
	CreationTime     string
	LastModifiedTime string
	Status           string
}

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
	if err != nil {
		return err
	}

	return nil
}

func GetBuckets(dir, name string) ([]byte, error) {
	file, err := os.Open(dir + "/" + name + ".csv")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if err != nil {
		return nil, err
	}

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	response := ListAllMyBucketsResult{Buckets: []Bucket{}}
	for _, line := range records[1:] {
		bucket := Bucket{
			Name:             line[0],
			CreationTime:     line[1],
			LastModifiedTime: line[2],
			Status:           line[3],
		}
		response.Buckets = append(response.Buckets, bucket)
	}

	xmlData, err := xml.MarshalIndent(response, "", " ")
	if err != nil {
		return nil, err
	}
	return xmlData, nil
}
