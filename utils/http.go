package utils

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"error": err.Error()})
}

func ParseQuery(r *http.Request, v any) error {
	err := decoder.Decode(v, r.URL.Query())
	if err != nil {
		return err
	}

	return nil
}

func ParseBody(r *http.Request, v any) error {
	if r.Body == nil {
		return fmt.Errorf("missing request body")
	}

	return json.NewDecoder(r.Body).Decode(v)
}
