package utils

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func DecodeJSON(r *http.Request, dst interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, "multipart/form-data") {
		return errors.New("expected application/json but received multipart/form-data. If you are uploading a file, please check the documentation for the correct S3 presigned URL upload flow")
	}
	return json.NewDecoder(r.Body).Decode(dst)
}
