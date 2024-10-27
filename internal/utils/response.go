package utils

import (
	"encoding/xml"
	"net/http"
)

type Error struct {
	Message  string `xml:"Message"`
	Resource string `xml:"Resource"`
}

type PutResult struct {
	Message string `xml:"Message"`
	Key     string `xml:"Key"`
}

type DeleteResult struct {
	Message string `xml:"Message"`
	Key     string `xml:"Key"`
}

type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

type Bucket struct {
	Name             string `xml:"Name"`
	CreationTime     string `xml:"CreationTime"`
	LastModifiedTime string `xml:"LastModifiedTime"`
	Status           string `xml:"Status"`
}

// XMLResponse marshals the given data into XML format and writes it to the response
func XMLResponse(w http.ResponseWriter, statusCode int, data interface{}) error {
	// Set the content-type
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(statusCode)

	// Marshal the data into XML format
	xmlData, err := xml.MarshalIndent(data, "", "   ")
	if err != nil {
		http.Error(w, "Failed to generate XML response", http.StatusInternalServerError)
		return err
	}

	// Write the XML header and the marshaled XML data
	w.Write([]byte(xml.Header))
	_, err = w.Write(xmlData)
	return err
}
