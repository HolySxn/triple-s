package bucket

import (
	"net/http"
	"os"
	"path"
	"time"

	"triple-s/internal/utils"
)

// CreateBucket creates a new bucket with the specified name and directory.
// It validates the bucket name, checks if the bucket already exists, and updates metadata if successful.
func CreateBucket(name string, dir string) (int, string) {
	// Validate bucket name
	if !utils.IsValidBucketName(name) {
		return http.StatusBadRequest, "Invalid Name"
	}

	// Check if bucket already exists in "buckets.csv"
	if flag, _, _ := utils.FindName(dir+"/buckets.csv", name); !flag {
		// Create directory for the new bucket
		bucket_dir := dir + "/" + name
		err := os.Mkdir(bucket_dir, os.ModePerm)
		if err != nil {
			return http.StatusInternalServerError, "Name is not Allowed"
		}

		// Prepare metadata for the bucket and set initial status as "InActive"
		metaData := []string{name, time.Now().Format("2006-01-02 15:04:05 MST"), time.Now().Format("2006-01-02 15:04:05 MST"), "InActive"}

		// Append metadata to "buckets.csv"
		err = utils.WriteCSV(dir+"/buckets.csv", metaData)
		if err != nil {
			return http.StatusInternalServerError, "Internal Server Error"
		}

		// Return success status and message
		return http.StatusOK, "Bucket was successfully created"
	} else {
		// Return conflict status if bucket already exists
		return http.StatusConflict, "Bucket Already Exists"
	}
}

// GetBucketsXML reads the bucket metadata file ("buckets.csv") and returns all buckets in XML format.
func GetBucketsXML(dir string) (utils.ListAllMyBucketsResult, int) {
	// Read all records from "buckets.csv"
	records, err := utils.ReadCSV(dir)
	if err != nil {
		return utils.ListAllMyBucketsResult{}, http.StatusInternalServerError
	}

	// Initialize response structure
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

	// Return bucket list and success status
	return response, http.StatusOK
}

// DeleteBucket deletes a bucket if it exists and is empty (status is "InActive").
// It removes the bucket directory and updates metadata in "buckets.csv".
func DeleteBucket(name string, dir string) (int, string) {
	// Check if the bucket exists and retrieve its metadata
	if flag, index, record := utils.FindName(dir+"/buckets.csv", name); flag {
		// Define paths for the bucket directory and metadata file
		bucket_dir := path.Join(dir, name)
		csv_dir := path.Join(dir, "buckets.csv")

		// Check if bucket is empty ("InActive" status)
		if record[3] == "InActive" {
			// Remove the bucket directory and its contents
			err := os.RemoveAll(bucket_dir)
			if err != nil {
				return http.StatusInternalServerError, "Internal Server Error"
			}

			// Remove the bucket record from "buckets.csv"
			err = utils.UpdateCSV(csv_dir, "delete", index, nil)
			if err != nil {
				return http.StatusInternalServerError, "Internal Server Error"
			}

			// Return no content status if deletion is successful
			return http.StatusNoContent, "Bucket was successfully deleted"
		} else {
			// Return conflict status if the bucket is not empty
			return http.StatusConflict, "Bucket is not empty"
		}
	} else {
		// Return not found status if the bucket does not exist
		return http.StatusNotFound, "Bucket is not found"
	}
}
