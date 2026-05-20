package api

import (
	"encoding/json"
	"net/http"
)

// APIResponse mirrors nestjs ApiResponse type shape.
type APIResponse struct {
	Success    bool   `json:"success"`
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
	Error      any    `json:"error,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeOK(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(w, status, APIResponse{
		Success:    true,
		StatusCode: status,
		Message:    message,
		Data:       data,
	})
}

func writeErr(w http.ResponseWriter, status int, message string, detail any) {
	writeJSON(w, status, APIResponse{
		Success:    false,
		StatusCode: status,
		Message:    message,
		Data:       nil,
		Error:      detail,
	})
}
