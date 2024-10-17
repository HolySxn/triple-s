package bucket

import (
	"encoding/csv"
	"encoding/xml"
	"net/http"
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

func GetBuckets(dir, name string) ([]byte, int) {
	file, err := os.Open(dir + "/" + name + ".csv")
	if err != nil {
		return nil, http.StatusInternalServerError
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
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
	bucket_dir := dir + "/" + name
	if !utils.IsExist(bucket_dir) {
		return http.StatusNotFound
	}

	files := utils.PahtFiles(bucket_dir)
	if len(files) > 1 {
		return http.StatusConflict
	}

	err := os.RemoveAll(bucket_dir)
	if err != nil{
		return http.StatusInternalServerError
	}

	return http.StatusNoContent
}
