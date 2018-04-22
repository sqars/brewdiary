package handlers

import (
	"encoding/json"
	"net/http"
)

type ingridientReq struct {
	Quantity     int  `json:"quantity"`
	IngridientID uint `json:"id"`
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
func decodeJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(&v)
}
