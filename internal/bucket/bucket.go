package bucket

import (
	"encoding/xml"
	"net/http"
	"os"
	"path"
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

func CreateBucket(name string, dir string) int {
	if !utils.IsValidBucketName(name) {
		return http.StatusBadRequest
	}

	bucket_dir := dir + "/" + name
	err := os.Mkdir(bucket_dir, os.ModePerm)
	if err != nil {
		return http.StatusConflict
	}

	metaData := []string{name, time.Now().Format("2006-01-02 15:04:05 MST"), time.Now().Format("2006-01-02 15:04:05 MST"), "Active"}
	err = utils.WriteCSV(dir+"/buckets.csv", metaData)
	if err != nil {
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

func GetBucketsXML(dir string) ([]byte, int) {
	records, err := utils.ReadCSV(dir)
	if err != nil {
		return nil, http.StatusInternalServerError
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
		return nil, http.StatusInternalServerError
	}
	return xmlData, http.StatusOK
}

func DeleteBucket(name string, dir string) int {
	if flag, index := utils.FindName(dir+"/buckets.csv", name); flag {
		bucket_dir := path.Join(dir, name)
		csv_dir := path.Join(bucket_dir, "objects.csv")
		if utils.IsEmptyCSV(csv_dir) {
			// Remove bucket
			err := os.RemoveAll(bucket_dir)
			if err != nil {
				return http.StatusInternalServerError
			}

			// Delete data from metadata
			err = utils.UpdateCSV(csv_dir, "delete", index, nil)
			if err != nil {
				return http.StatusInternalServerError
			}

			return http.StatusNoContent
		} else {
			return http.StatusConflict
		}
	} else {
		return http.StatusNotFound
	}
}
