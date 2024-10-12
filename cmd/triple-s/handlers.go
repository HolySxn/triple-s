package main

import (
	"fmt"
	"net/http"
)

func handlerMangaer(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPut:
		createBucket(w, r)
	case http.MethodGet:
		getBuckets(w, r)
	case http.MethodDelete:
		deleteBucket(w, r)
	default:
		http.Error(w, "Method not Allowed", http.StatusMethodNotAllowed)
		return
	}
}

func createBucket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Create new bucket %v\n", r.URL.Path[1:])
}

func getBuckets(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Get all buckets\n")
}

func deleteBucket(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Delete bucket %v\n", r.URL.Path[1:])
}
