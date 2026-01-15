package adapterhttp

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, data any, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(Response{
		Success: true,
		Message: msg,
		Data:    data,
	})
	if err != nil {
		return
	}
}

func WriteErrorJSON(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(Response{
		Success: false,
		Message: msg,
		Data:    nil,
	})
	if err != nil {
		return
	}
}
