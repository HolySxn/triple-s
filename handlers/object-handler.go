package handlers

import (
	"io"
	"net/http"
	"path"
	"time"

	"triple-s/internal/object"
	"triple-s/internal/utils"
)

func ObjectHnadler(dir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bucket_name := r.PathValue("BucketName")
		if f, _ := utils.FindName(path.Join(dir, "buckets.csv"), bucket_name); !f {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		object_name := r.PathValue("ObjectName")
		switch r.Method {
		case http.MethodPut:
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
			record := []string{object_name, http.DetectContentType(data[:512]), time.Now().Format("2006-01-02 15:04:05 MST")}
				if f, index := utils.FindName(csv_dir, object_name); f {
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

			w.WriteHeader(status)
			w.Write([]byte("Object was added successfully!"))
		}
	}
}
