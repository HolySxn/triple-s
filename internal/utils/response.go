package utils

import (
	"encoding/xml"
	"net/http"
)

// Error structure represents an XML response for error messages.
// It includes a message and the resource that caused the error.
type Error struct {
	Code     string `xml:"Code"`
	Message  string `xml:"Message"`
	Resource string `xml:"Resource"`
}

// PutResult structure represents a successful response for PUT requests,
// containing a message and the object key (identifier).
type PutResult struct {
	Message string `xml:"Message"`
	Key     string `xml:"Key"`
}

// DeleteResult structure represents a successful response for DELETE requests,
// containing a message and the key of the deleted object.
type DeleteResult struct {
	Message string `xml:"Message"`
	Key     string `xml:"Key"`
}

// ListAllMyBucketsResult structure represents the response for listing all buckets.
// It contains a slice of Bucket structs nested within the <Buckets> XML element.
type ListAllMyBucketsResult struct {
	Buckets []Bucket `xml:"Buckets>Bucket"`
}

// Bucket structure represents a single bucket's metadata in XML format,
// including name, creation time, last modified time, and status.
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
