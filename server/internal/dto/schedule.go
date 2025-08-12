package dto

import "time"

type ScheduleCreateDTO struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartTime   string  `json:"start_time"`             // ISO 8601 format
	EndTime     string  `json:"end_time"`               // ISO 8601 format
	RepeatType  string  `json:"repeat_type"`            // e.g. "none", "daily", "weekly"
	RepeatUntil *string `json:"repeat_until,omitempty"` // ISO 8601 format, optional
}

type ScheduleUpdateDTO struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	StartTime   *string `json:"start_time,omitempty"`
	EndTime     *string `json:"end_time,omitempty"`
	IsDone      *bool   `json:"is_done,omitempty"`
	RepeatType  *string `json:"repeat_type,omitempty"`
	RepeatUntil *string `json:"repeat_until,omitempty"`
}

type ScheduleResponseDTO struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	IsDone      bool    `json:"is_done"`
	CreatedAt   string  `json:"created_at"`
	RepeatType  string  `json:"repeat_type"`
	RepeatUntil *string `json:"repeat_until,omitempty"`
}

type CreateScheduleRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartTime   time.Time  `json:"start_time"`
	EndTime     time.Time  `json:"end_time"`
	RepeatType  string     `json:"repeat_type"`
	RepeatUntil *time.Time `json:"repeat_until,omitempty"`
}

type UpdateScheduleRequest struct {
	Title       *string    `json:"title,omitempty"`
	Description *string    `json:"description,omitempty"`
	StartTime   *time.Time `json:"start_time,omitempty"`
	EndTime     *time.Time `json:"end_time,omitempty"`
	IsDone      *bool      `json:"is_done,omitempty"`
	RepeatType  *string    `json:"repeat_type,omitempty"`
	RepeatUntil *time.Time `json:"repeat_until,omitempty"`
}

// PaginatedResponse is a generic wrapper for paginated API responses
type PaginatedResponse struct {
	Data       interface{} `json:"data"`        // The actual list of items
	Page       int         `json:"page"`        // Current page number
	Limit      int         `json:"limit"`       // Items per page
	TotalItems int         `json:"total_items"` // Total number of items in DB
	TotalPages int         `json:"total_pages"` // Total number of pages
}
