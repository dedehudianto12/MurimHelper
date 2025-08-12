package dto

import "time"

type ScheduleFilter struct {
	IsDone      *bool
	RepeatType  string
	Search      string
	StartAfter  *time.Time
	StartBefore *time.Time
	SortBy      string
	SortOrder   string
}
