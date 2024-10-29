package handlers

import (
	"net/http"
	"strings"
	"triple-s/internal/bucket"
	"triple-s/internal/utils"
)

// BucketHandler handles HTTP requests related to bucket operations.
// It takes a directory path as an argument and returns an http.HandlerFunc to manage
// different HTTP methods (PUT, GET, DELETE) for bucket-related actions.
func BucketHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			// Extract bucket name from URL path
			bucketName := strings.TrimPrefix(r.URL.Path, "/")

			// Create a new bucket with the specified name
			status, message := bucket.CreateBucket(bucketName, dir)
			if status != http.StatusOK {
				utils.XMLResponse(w, status,
					utils.Error{
						Code:     http.StatusText(status),
						Message:  message,
						Resource: r.URL.Path,
					})
				return
			}

			// Create a CSV file for storing objects metadata in the new bucket's directory
			err := utils.CreateCSV(dir + "/" + bucketName + "/objects.csv")
			if err != nil {
				utils.XMLResponse(w, status,
					utils.Error{
						Code:     http.StatusText(http.StatusInternalServerError),
						Message:  "Can not create metadata",
						Resource: r.URL.Path,
					})
				return
			}

			// Write header columns for the objects CSV file
			err = utils.WriteCSV(dir+"/"+bucketName+"/objects.csv", []string{"ObjectKey", "Size", "ContentType", "LastModified"})
			if err != nil {
				utils.XMLResponse(w, status,
					utils.Error{
						Code:     http.StatusText(http.StatusInternalServerError),
						Message:  "Can not add new record into objects metadata",
						Resource: r.URL.Path,
					})
				return
			}

			// Successfully created bucket and CSV file, respond with success in XML format
			utils.XMLResponse(w, status, utils.PutResult{Message: message, Key: bucketName})
		case http.MethodGet:
			// Handle GET requests to list all buckets; only the root path ("/") is allowed
			if r.URL.Path != "/" {
				utils.XMLResponse(w, http.StatusBadRequest,
					utils.Error{
						Code:     http.StatusText(http.StatusBadRequest),
						Message:  "Only the root path / is allowed",
						Resource: r.URL.Path,
					})
				return
			}

			// Retrieve bucket list in XML format from CSV file
			data, status := bucket.GetBucketsXML(dir + "/buckets.csv")
			if status != http.StatusOK {
				utils.XMLResponse(w, status,
					utils.Error{
						Code:     http.StatusText(http.StatusInternalServerError),
						Message:  "Can not read buckets metadata",
						Resource: r.URL.Path,
					})
				return
			}

			// Respond with the bucket list in XML format
			utils.XMLResponse(w, status, data)
		case http.MethodDelete:
			// Extract bucket name from URL path for deletion
			bucketName := strings.TrimPrefix(r.URL.Path, "/")

			// Attempt to delete the specified bucket
			status, message := bucket.DeleteBucket(bucketName, dir)
			if status != http.StatusNoContent {
				utils.XMLResponse(w, status,
					utils.Error{
						Code:     http.StatusText(status),
						Message:  message,
						Resource: r.URL.Path,
					})
				return
			}

			// Successful deletion; respond with `204 No Content` status code without body
			w.WriteHeader(status)
		default:
			// If the HTTP method is not PUT, GET, or DELETE, respond with a 405 Method Not Allowed
			utils.XMLResponse(w, http.StatusMethodNotAllowed,
				utils.Error{
					Code:     http.StatusText(http.StatusMethodNotAllowed),
					Message:  "Method" + r.Method + "is not allowed",
					Resource: r.Method,
				})
		}
	}
}
