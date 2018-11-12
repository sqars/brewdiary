package handlers

import (
	"encoding/json"
	"net/http"
)

type ok interface {
	OK() error
}

// ResponseJSON returns response as JSON with passed payload
func responseJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}

// DecodeJSON decodes JSON request
func decodeJSON(r *http.Request, v ok, validate bool) error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&v)
	if err != nil {
		return err
	}
	if validate {
		return v.OK()
	}
	return nil
}
