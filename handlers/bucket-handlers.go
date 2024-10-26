package handlers

import (
	"encoding/xml"
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
			status := bucket.CreateBucket(bucketName, dir)
			if status != http.StatusOK {
				http.Error(w, http.StatusText(status), status)
				return
			}

			err := utils.CreateCSV(dir + "/" + bucketName + "/objects.csv")
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			err = utils.WriteCSV(dir+"/"+bucketName+"/objects.csv", []string{"ObjectKey", "Size", "ContentType", "LastModified"})
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(status)
			w.Write([]byte("Bucket created successfully"))
		case http.MethodGet:
			if r.URL.Path != "/" {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}
			xmlData, status := bucket.GetBucketsXML(dir + "/buckets.csv")
			if status != http.StatusOK {
				http.Error(w, http.StatusText(status), status)
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(status)
			w.Write([]byte(xml.Header))
			w.Write(xmlData)
		case http.MethodDelete:
			bucketName := strings.TrimPrefix(r.URL.Path, "/")
			status := bucket.DeleteBucket(bucketName, dir)
			if status != http.StatusNoContent {
				http.Error(w, http.StatusText(status), status)
				return
			}
			w.WriteHeader(status)
			w.Write([]byte("Bucket was successully deleted"))
		default:
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		}
	}
}
