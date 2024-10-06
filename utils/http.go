package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func WriteResponseJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteResponseError(w http.ResponseWriter, status int, err error) {
	WriteResponseJSON(w, status, map[string]any{
		"error": err.Error(),
	})
}

func WriteResponseData(w http.ResponseWriter, v any) {
	WriteResponseJSON(w, http.StatusOK, v)
}

func ParseQuery(r *http.Request, obj any) error {
	return decoder.Decode(obj, r.URL.Query())
}

func ParseBody(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}
