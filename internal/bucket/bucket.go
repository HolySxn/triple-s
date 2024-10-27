package bucket

import (
	"net/http"
	"os"
	"path"
	"time"

	"triple-s/internal/utils"
)

func CreateBucket(name string, dir string) (int, string) {
	if !utils.IsValidBucketName(name) {
		return http.StatusBadRequest, "Invalid Name"
	}

	if flag, _, _ := utils.FindName(dir+"/buckets.csv", name); !flag {
		bucket_dir := dir + "/" + name
		err := os.Mkdir(bucket_dir, os.ModePerm)
		if err != nil {
			return http.StatusInternalServerError, "Name is not Allowed"
		}

		metaData := []string{name, time.Now().Format("2006-01-02 15:04:05 MST"), time.Now().Format("2006-01-02 15:04:05 MST"), "InActive"}
		err = utils.WriteCSV(dir+"/buckets.csv", metaData)
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}

		return http.StatusOK, "Bucket was successfully created"
	} else {
		return http.StatusConflict, "Bucket Already Exists"
	}
}

func GetBucketsXML(dir string) (utils.ListAllMyBucketsResult, int) {
	records, err := utils.ReadCSV(dir)
	if err != nil {
		return utils.ListAllMyBucketsResult{}, http.StatusInternalServerError
	}

	response := utils.ListAllMyBucketsResult{Buckets: []utils.Bucket{}}
	for _, line := range records[1:] {
		bucket := utils.Bucket{
			Name:             line[0],
			CreationTime:     line[1],
			LastModifiedTime: line[2],
			Status:           line[3],
		}
		response.Buckets = append(response.Buckets, bucket)
	}

	return response, http.StatusOK
}

func DeleteBucket(name string, dir string) (int, string) {
	if flag, index, record := utils.FindName(dir+"/buckets.csv", name); flag {
		bucket_dir := path.Join(dir, name)
		csv_dir := path.Join(dir, "buckets.csv")
		if record[3] == "InActive" {
			// Remove bucket
			err := os.RemoveAll(bucket_dir)
			if err != nil {
				return http.StatusInternalServerError, "Internal Server Error"
			}

			// Delete data from metadata
			err = utils.UpdateCSV(csv_dir, "delete", index, nil)
			if err != nil {
				return http.StatusInternalServerError, "Internal Server Error"
			}

			return http.StatusNoContent, "Bucket was successfully deleted"
		} else {
			return http.StatusConflict, "Bucket is not empty"
		}
	} else {
		return http.StatusNotFound, "Bucket is not found"
	}
}
