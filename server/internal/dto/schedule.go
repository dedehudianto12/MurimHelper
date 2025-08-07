package dto

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
