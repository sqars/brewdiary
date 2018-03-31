package utils

import (
	"encoding/json"
	"net/http"
)

// ResponseJSON returns response as JSON with passed payload
func ResponseJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload != nil {
		json.NewEncoder(w).Encode(payload)
	}
}
