package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"triple-s/internal/bucket"
	"triple-s/internal/utils"
)

func BucketHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			bucketName := r.URL.Query().Get("BucketName")
			status := bucket.CreateBucket(bucketName, dir)
			if status != http.StatusOK{
				http.Error(w, http.StatusText(status), status)
				return
			}
			utils.CreateCSV(dir+"/"+bucketName, "objects", []string{"ObjectKey", "Size", "ContentType", "LastModified"})
			w.WriteHeader(status)
			w.Write([]byte("Bucket created successfully"))
		case http.MethodGet:
			xmlData, status := bucket.GetBuckets(dir, "buckets")
			if status != http.StatusOK{
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
			if status != http.StatusNoContent{
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

func ObjectHnadler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Object handler %v\n", r.URL.Path)
}
