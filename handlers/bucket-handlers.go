package handlers

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"triple-s/internal/bucket"
)

func BucketHandler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			bucketName := strings.TrimPrefix(r.URL.Path, "/")
			err := bucket.CreateBucket(bucketName, dir)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusCreated)
			w.Write([]byte("Bucket created successfully"))
			return
		case http.MethodGet:
			xmlData, err := bucket.GetBuckets(dir, "buckets")
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("Content-Type", "application/xml")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(xml.Header))
			w.Write(xmlData)
		case http.MethodDelete:

		default:
			http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
			return
		}
	}
}

func ObjectHnadler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Object handler %v\n", r.URL.Path)
}
