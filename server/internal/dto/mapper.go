package dto

import (
	"errors"
	"murim-helper/internal/domain"
	"strings"
	"time"
)

// ======================
// Response DTOs
// ======================

func ToScheduleResponseDTO(s domain.Schedule) ScheduleResponseDTO {
	var repeatUntil *string
	if s.RepeatUntil != nil {
		str := s.RepeatUntil.Format(time.RFC3339)
		repeatUntil = &str
	}

	return ScheduleResponseDTO{
		ID:          s.ID,
		Title:       s.Title,
		Description: s.Description,
		StartTime:   s.StartTime.Format(time.RFC3339),
		EndTime:     s.EndTime.Format(time.RFC3339),
		IsDone:      s.IsDone,
		CreatedAt:   s.CreatedAt.Format(time.RFC3339),
		RepeatType:  s.RepeatType,
		RepeatUntil: repeatUntil,
	}
}

func ToScheduleResponseDTOs(schedules []domain.Schedule) []ScheduleResponseDTO {
	result := make([]ScheduleResponseDTO, len(schedules))
	for i, s := range schedules {
		result[i] = ToScheduleResponseDTO(s)
	}
	return result
}

// ======================
// Request DTOs
// ======================

type GenerateScheduleRequest struct {
	Description string `json:"description"`
}

func (r GenerateScheduleRequest) Validate() error {
	if strings.TrimSpace(r.Description) == "" {
		return errors.New("description is required")
	}
	return nil
}

func (r *CreateScheduleRequest) Validate() error {
	if strings.TrimSpace(r.Title) == "" {
		return errors.New("title is required")
	}
	if r.StartTime.IsZero() || r.EndTime.IsZero() {
		return errors.New("start_time and end_time are required")
	}
	if r.StartTime.After(r.EndTime) {
		return errors.New("start_time must be before end_time")
	}
	if r.RepeatType == "" {
		r.RepeatType = "none"
	}
	return nil
}

func (r CreateScheduleRequest) ToDomain() domain.Schedule {
	return domain.Schedule{
		Title:       r.Title,
		Description: r.Description,
		StartTime:   r.StartTime,
		EndTime:     r.EndTime,
		IsDone:      false,
		RepeatType:  r.RepeatType,
		RepeatUntil: r.RepeatUntil,
	}
}

func (r UpdateScheduleRequest) Validate() error {
	if r.StartTime != nil && r.EndTime != nil && r.StartTime.After(*r.EndTime) {
		return errors.New("start_time must be before end_time")
	}
	return nil
}

func (r UpdateScheduleRequest) ToDomain(existing domain.Schedule) domain.Schedule {
	if r.Title != nil {
		existing.Title = *r.Title
	}
	if r.Description != nil {
		existing.Description = *r.Description
	}
	if r.StartTime != nil {
		existing.StartTime = *r.StartTime
	}
	if r.EndTime != nil {
		existing.EndTime = *r.EndTime
	}
	if r.IsDone != nil {
		existing.IsDone = *r.IsDone
	}
	if r.RepeatType != nil {
		existing.RepeatType = *r.RepeatType
	}
	if r.RepeatUntil != nil {
		existing.RepeatUntil = r.RepeatUntil
	}
	return existing
}
