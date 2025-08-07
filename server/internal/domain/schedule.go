package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Schedule struct {
	ID          string    `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description string    `db:"description" json:"description"`
	StartTime   time.Time `db:"start_time" json:"start_time"`
	EndTime     time.Time `db:"end_time" json:"end_time"`
	IsDone      bool      `db:"is_done" json:"is_done"`
	CreatedAt   string    `db:"created_at" json:"created_at"`
}

func ParseSchedulesFromJSON(jsonStr string) ([]Schedule, error) {
	var rawItems []struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		StartTime   string `json:"start_time"` // ISO format
		EndTime     string `json:"end_time"`
	}

	err := json.Unmarshal([]byte(jsonStr), &rawItems)
	if err != nil {
		return nil, err
	}

	var schedules []Schedule
	for _, item := range rawItems {
		start, err1 := time.Parse(time.RFC3339, item.StartTime)
		end, err2 := time.Parse(time.RFC3339, item.EndTime)
		if err1 != nil || err2 != nil {
			fmt.Printf("Error parsing times:\n  start: %v\n  end: %v\n", err1, err2)
			continue
		}

		schedules = append(schedules, Schedule{
			ID:          uuid.NewString(),
			Title:       item.Title,
			Description: item.Description,
			StartTime:   start,
			EndTime:     end,
			IsDone:      false,
		})
	}

	return schedules, nil
}
