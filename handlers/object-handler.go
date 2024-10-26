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

func ObjectHnadler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket_name := r.PathValue("BucketName")
		if f, _, _ := utils.FindName(path.Join(dir, "buckets.csv"), bucket_name); !f {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		object_name := r.PathValue("ObjectName")
		switch r.Method {
		case http.MethodPut:
			if !utils.IsValidObjectName(object_name) {
				http.Error(w, http.StatusText(http.StatusConflict), http.StatusConflict)
				return
			}
			ObjectPut(w, r, dir, bucket_name, object_name)
		case http.MethodGet:
			csv_dir := path.Join(dir, bucket_name, "objects.csv")
			if f, _, _ := utils.FindName(csv_dir, object_name); f {
				data, status := object.GetObject(path.Join(dir, bucket_name, object_name))
				if status != http.StatusOK {
					http.Error(w, http.StatusText(status), status)
					return
				}

				w.WriteHeader(status)
				w.Write(data)
			} else {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
		case http.MethodDelete:
			status := object.DeleteObject(object_name, path.Join(dir, bucket_name))
			if status != http.StatusOK {
				http.Error(w, http.StatusText(status), status)
				return
			}

			bucket_dir := path.Join(dir, "buckets.csv")
			_, index, record := utils.FindName(bucket_dir, bucket_name)
			record[2] = time.Now().Format("2006-01-02 15:04:05 MST")
			if utils.IsEmptyCSV(path.Join(dir, bucket_name, "objects.csv")) {
				record[3] = "InActive"
			}
			err := utils.UpdateCSV(bucket_dir, "update", index, record)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			w.WriteHeader(status)
			w.Write([]byte("Object was deleted successfully!"))

		}
	}
}

func ObjectPut(w http.ResponseWriter, r *http.Request, dir, bucket_name, object_name string) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	object_dir := path.Join(dir, bucket_name, object_name)
	status := object.PutObject(data, object_dir)
	if status != http.StatusOK {
		http.Error(w, http.StatusText(status), status)
		return
	}

	csv_dir := path.Join(dir, bucket_name, "objects.csv")
	record := []string{
		object_name,
		strconv.Itoa(int(r.ContentLength)),
		r.Header.Get("Content-Type"),
		time.Now().Format("2006-01-02 15:04:05 MST"),
	}
	if f, index, _ := utils.FindName(csv_dir, object_name); f {
		err = utils.UpdateCSV(csv_dir, "update", index, record)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	} else {
		err = utils.WriteCSV(csv_dir, record)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

	bucket_dir := path.Join(dir, "buckets.csv")
	_, index, record := utils.FindName(bucket_dir, bucket_name)
	record[2] = time.Now().Format("2006-01-02 15:04:05 MST")
	record[3] = "Active"
	err = utils.UpdateCSV(bucket_dir, "update", index, record)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write([]byte("Object was added successfully!"))
}
