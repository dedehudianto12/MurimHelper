package dto

import (
	"murim-helper/internal/domain"
	"time"
)

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
