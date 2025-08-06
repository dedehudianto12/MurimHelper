package httphelper

import (
	"encoding/json"
	"net/http"
	"time"

	"murim-helper/internal/model"
)

// Success returns a well-structured successful response
func Success(w http.ResponseWriter, r *http.Request, statusCode int, message string, payload interface{}) {
	resp := model.ApiResponse{
		Status:    "success",
		Message:   message,
		Payload:   payload,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      r.URL.Path,
	}

	writeJSON(w, statusCode, resp)
}

// Error returns a well-structured error response
func Error(w http.ResponseWriter, r *http.Request, statusCode int, message string, code int) {
	resp := model.ApiResponse{
		Status:    "error",
		Message:   message,
		Code:      code,
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Path:      r.URL.Path,
	}

	writeJSON(w, statusCode, resp)
}

// Helper to write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}
