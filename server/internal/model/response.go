package model

type ApiResponse struct {
	Status    string      `json:"status"`              // "success" or "error"
	Message   string      `json:"message,omitempty"`   // Human-readable message
	Payload   interface{} `json:"payload,omitempty"`   // Actual response data
	Timestamp string      `json:"timestamp,omitempty"` // RFC3339 timestamp
	Path      string      `json:"path,omitempty"`      // Request URL path
	Code      int         `json:"code,omitempty"`      // Optional internal code
	Meta      interface{} `json:"meta,omitempty"`      // For paginated responses etc.
}
