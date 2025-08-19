package api

import (
	"encoding/json"
	"net/http"
)

func OKResponse(w http.ResponseWriter, data any) {
	// Return the products as a JSON response
	w.Header().Set("Content-Type", "application/json")

	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	http.Error(w, message, status)
}
