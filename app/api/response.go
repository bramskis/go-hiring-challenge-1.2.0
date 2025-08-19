package api

import (
	"encoding/json"
	"net/http"
)

type ErrorJSON struct {
	Error string `json:"error"`
}

func OKResponse(w http.ResponseWriter, data any) {
	writeData(w, http.StatusOK, data)
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	writeData(w, status, ErrorJSON{Error: message})
}

func writeData(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
