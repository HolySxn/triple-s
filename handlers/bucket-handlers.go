package handlers

import (
	"encoding/csv"
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"

	"triple-s/internal/bucket"
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
			GetBuckets(dir, "buckets", w, r)
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

func GetBuckets(dir, name string, w http.ResponseWriter, r *http.Request) {
	file, err := os.Open(dir + "/" + name + ".csv")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := ListAllMyBucketsResult{Buckets: []Bucket{}}
	for i := 1; i < len(records); i++{
		bucket := Bucket{
			Name:             records[i][0],
			CreationTime:     records[i][1],
			LastModifiedTime: records[i][2],
			Status:           records[i][3],
		}
		response.Buckets = append(response.Buckets, bucket)
	}

	x, err := xml.MarshalIndent(response, "", " ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(xml.Header))
	w.Write(x)
}
