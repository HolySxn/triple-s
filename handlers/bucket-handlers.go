package handlers

import (
	"net/http"
	"strings"

	"triple-s/internal/bucket"
	"triple-s/internal/utils"
)

func BucketHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			bucketName := strings.TrimPrefix(r.URL.Path, "/")
			status, message := bucket.CreateBucket(bucketName, dir)
			if status != http.StatusOK {
				utils.XMLResponse(w, status, utils.Error{Message: message, Resource: r.URL.Path})
				return
			}

			err := utils.CreateCSV(dir + "/" + bucketName + "/objects.csv")
			if err != nil {
				utils.XMLResponse(w, status, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
				return
			}

			err = utils.WriteCSV(dir+"/"+bucketName+"/objects.csv", []string{"ObjectKey", "Size", "ContentType", "LastModified"})
			if err != nil {
				utils.XMLResponse(w, status, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
				return
			}

			utils.XMLResponse(w, status, utils.PutResult{Message: message, Key: bucketName})
		case http.MethodGet:
			if r.URL.Path != "/" {
				utils.XMLResponse(w, http.StatusBadRequest, utils.Error{Message: "Bad Request", Resource: r.URL.Path})
				return
			}
			data, status := bucket.GetBucketsXML(dir + "/buckets.csv")
			if status != http.StatusOK {
				utils.XMLResponse(w, status, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
				return
			}

			utils.XMLResponse(w, status, data)
		case http.MethodDelete:
			bucketName := strings.TrimPrefix(r.URL.Path, "/")
			status, message := bucket.DeleteBucket(bucketName, dir)
			if status != http.StatusNoContent {
				utils.XMLResponse(w, status, utils.Error{Message: message, Resource: r.URL.Path})
				return
			}

			w.WriteHeader(status)
		default:
			utils.XMLResponse(w, http.StatusMethodNotAllowed, utils.Error{Message: "Method is not allowed", Resource: r.Method})
		}
	}
}
