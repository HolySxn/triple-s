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
			utils.XMLResponse(w, http.StatusNotFound, utils.Error{Message: "Bucket is not found", Resource: r.URL.Path})
			return
		}

		object_name := r.PathValue("ObjectName")
		switch r.Method {
		case http.MethodPut:
			if !utils.IsValidObjectName(object_name) {
				utils.XMLResponse(w, http.StatusBadRequest, utils.Error{Message: "Invalid object name", Resource: r.URL.Path})
				return
			}
			ObjectPut(w, r, dir, bucket_name, object_name)
		case http.MethodGet:
			csv_dir := path.Join(dir, bucket_name, "objects.csv")
			if f, _, _ := utils.FindName(csv_dir, object_name); f {
				data, status := object.GetObject(path.Join(dir, bucket_name, object_name))
				if status != http.StatusOK {
					utils.XMLResponse(w, http.StatusInternalServerError, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
					return
				}

				w.WriteHeader(status)
				w.Write(data)
			} else {
				utils.XMLResponse(w, http.StatusNotFound, utils.Error{Message: "Object is not found", Resource: r.URL.Path})
				return
			}
		case http.MethodDelete:
			status, message := object.DeleteObject(object_name, path.Join(dir, bucket_name))
			if status != http.StatusOK {
				utils.XMLResponse(w, http.StatusInternalServerError, utils.Error{Message: message, Resource: r.URL.Path})
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
				utils.XMLResponse(w, http.StatusInternalServerError, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
				return
			}

			utils.XMLResponse(w, status, utils.DeleteResult{Message: message, Key: object_name})
		default:
			utils.XMLResponse(w, http.StatusMethodNotAllowed, utils.Error{Message: "Method is not allowed", Resource: r.Method})
		}
	}
}

func ObjectPut(w http.ResponseWriter, r *http.Request, dir, bucket_name, object_name string) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		utils.XMLResponse(w, http.StatusInternalServerError, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
		return
	}

	object_dir := path.Join(dir, bucket_name, object_name)
	status, message := object.PutObject(data, object_dir)
	if status != http.StatusOK {
		utils.XMLResponse(w, status, utils.Error{Message: message, Resource: r.URL.Path})
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
			utils.XMLResponse(w, http.StatusInternalServerError, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
			return
		}
	} else {
		err = utils.WriteCSV(csv_dir, record)
		if err != nil {
			utils.XMLResponse(w, http.StatusInternalServerError, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
			return
		}
	}

	bucket_dir := path.Join(dir, "buckets.csv")
	_, index, record := utils.FindName(bucket_dir, bucket_name)
	record[2] = time.Now().Format("2006-01-02 15:04:05 MST")
	if record[3] == "InActive" {
		record[3] = "Active"
	}
	err = utils.UpdateCSV(bucket_dir, "update", index, record)
	if err != nil {
		utils.XMLResponse(w, http.StatusInternalServerError, utils.Error{Message: "Internal Server Error", Resource: r.URL.Path})
		return
	}

	utils.XMLResponse(w, status, utils.PutResult{Message: message, Key: object_name})
}
