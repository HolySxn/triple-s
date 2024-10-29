package handlers

import (
	"io"
	"net/http"
	"path"
	"strconv"
	"time"
	"triple-s/internal/object"
	"triple-s/internal/utils"
)

// ObjectHandler handles requests for object-related operations in a bucket.
// It takes a directory path as an argument and returns an http.HandlerFunc
// to manage different HTTP methods (PUT, GET, DELETE) for object-related actions.
func ObjectHnadler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract bucket name from the URL path
		bucket_name := r.PathValue("BucketName")

		// Verify if the bucket exists by checking "buckets.csv"
		if f, _, _ := utils.FindName(path.Join(dir, "buckets.csv"), bucket_name); !f {
			utils.XMLResponse(w, http.StatusNotFound,
				utils.Error{
					Code:     http.StatusText(http.StatusNotFound),
					Message:  "Bucket is not found",
					Resource: r.URL.Path,
				})
			return
		}

		// Extract object name from the URL path
		object_name := r.PathValue("ObjectName")
		switch r.Method {
		case http.MethodPut:
			// Validate object name
			if !utils.IsValidObjectName(object_name) {
				utils.XMLResponse(w, http.StatusBadRequest,
					utils.Error{
						Code:     http.StatusText(http.StatusBadRequest),
						Message:  "Invalid object name",
						Resource: r.URL.Path,
					})
				return
			}

			// Call ObjectPut to handle object upload
			ObjectPut(w, r, dir, bucket_name, object_name)
		case http.MethodGet:
			// Check if object exists by searching in "objects.csv"
			csv_dir := path.Join(dir, bucket_name, "objects.csv")
			if f, _, record := utils.FindName(csv_dir, object_name); f {
				// Retrieve object data
				data, status := object.GetObject(path.Join(dir, bucket_name, object_name))
				if status != http.StatusOK {
					utils.XMLResponse(w, http.StatusInternalServerError,
						utils.Error{
							Code:     http.StatusText(http.StatusInternalServerError),
							Message:  "Can not get object",
							Resource: r.URL.Path,
						})
					return
				}

				// Write object data to response
				w.Header().Set("Key", record[0])
				w.Header().Set("Content-Legth", strconv.Itoa(len(data)))
				w.Header().Set("Last-Modified", record[3])
				w.WriteHeader(status)
				w.Write(data)
			} else {
				// If object is not found, respond with 404 Not Found
				utils.XMLResponse(w, http.StatusNotFound,
					utils.Error{
						Code:     http.StatusText(http.StatusNotFound),
						Message:  "Object is not found",
						Resource: r.URL.Path,
					})
				return
			}
		case http.MethodDelete:
			// Attempt to delete the object
			status, message := object.DeleteObject(object_name, path.Join(dir, bucket_name))
			if status != http.StatusOK {
				utils.XMLResponse(w, status,
					utils.Error{
						Code:     http.StatusText(status),
						Message:  message,
						Resource: r.URL.Path,
					})
				return
			}

			// Update bucket's metadata in "buckets.csv" upon deletion
			bucket_dir := path.Join(dir, "buckets.csv")
			_, index, record := utils.FindName(bucket_dir, bucket_name)
			record[2] = time.Now().Format("2006-01-02 15:04:05 MST")
			if utils.IsEmptyCSV(path.Join(dir, bucket_name, "objects.csv")) {
				record[3] = "InActive"
			}
			err := utils.UpdateCSV(bucket_dir, "update", index, record)
			if err != nil {
				utils.XMLResponse(w, http.StatusInternalServerError,
					utils.Error{
						Code:     http.StatusText(http.StatusInternalServerError),
						Message:  "Can not update buckets metadata",
						Resource: r.URL.Path,
					})
				return
			}

			// Respond with a confirmation message in XML format
			w.Header().Set("Connetction", "Close")
			w.Header().Set("Content-Lenght", "0")
			utils.XMLResponse(w, status, utils.DeleteResult{Message: message, Key: object_name})
		default:
			// Respond with 405 Method Not Allowed for unsupported HTTP methods
			utils.XMLResponse(w, http.StatusMethodNotAllowed,
				utils.Error{
					Code:     http.StatusText(http.StatusMethodNotAllowed),
					Message:  "Method" + r.Method + "is not allowed",
					Resource: r.Method,
				})
		}
	}
}

// ObjectPut handles the uploading of an object to a specified bucket.
// It writes the object data to the appropriate location and updates metadata in CSV files.
func ObjectPut(w http.ResponseWriter, r *http.Request, dir, bucket_name, object_name string) {
	// Read the body of the PUT request to get the object data
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.XMLResponse(w, http.StatusInternalServerError,
			utils.Error{
				Code:     http.StatusText(http.StatusInternalServerError),
				Message:  "Can not get object data",
				Resource: r.URL.Path,
			})
		return
	}

	// Define the path where the object data will be stored
	object_dir := path.Join(dir, bucket_name, object_name)

	// Store the object data using object.PutObject
	status, message := object.PutObject(data, object_dir)
	if status != http.StatusOK {
		utils.XMLResponse(w, status,
			utils.Error{
				Code:     http.StatusText(status),
				Message:  message,
				Resource: r.URL.Path,
			})
		return
	}

	// Record object metadata in "objects.csv" for the specified bucket
	csv_dir := path.Join(dir, bucket_name, "objects.csv")
	record := []string{
		object_name,
		strconv.Itoa(int(r.ContentLength)),
		r.Header.Get("Content-Type"),
		time.Now().Format("2006-01-02 15:04:05 MST"),
	}

	// Check if the object already exists in "objects.csv" to update or add
	if f, index, _ := utils.FindName(csv_dir, object_name); f {
		// Update existing object entry
		err = utils.UpdateCSV(csv_dir, "update", index, record)
		if err != nil {
			utils.XMLResponse(w, http.StatusInternalServerError,
				utils.Error{
					Code:     http.StatusText(http.StatusInternalServerError),
					Message:  "Can not update objects metadata",
					Resource: r.URL.Path,
				})
			return
		}
	} else {
		// Write new object entry to CSV
		err = utils.WriteCSV(csv_dir, record)
		if err != nil {
			utils.XMLResponse(w, http.StatusInternalServerError,
				utils.Error{
					Code:     http.StatusText(http.StatusInternalServerError),
					Message:  "Can not add new record into objects metadata",
					Resource: r.URL.Path,
				})
			return
		}
	}

	// Update bucket metadata in "buckets.csv" after object creation
	bucket_dir := path.Join(dir, "buckets.csv")
	_, index, record := utils.FindName(bucket_dir, bucket_name)
	record[2] = time.Now().Format("2006-01-02 15:04:05 MST")
	if record[3] == "InActive" {
		record[3] = "Active"
	}
	err = utils.UpdateCSV(bucket_dir, "update", index, record)
	if err != nil {
		utils.XMLResponse(w, http.StatusInternalServerError,
			utils.Error{
				Code:     http.StatusText(http.StatusInternalServerError),
				Message:  "Can not update buckets metadata",
				Resource: r.URL.Path,
			})
		return
	}

	// Respond with a confirmation message in XML format
	utils.XMLResponse(w, status, utils.PutResult{Message: message, Key: object_name})
}
