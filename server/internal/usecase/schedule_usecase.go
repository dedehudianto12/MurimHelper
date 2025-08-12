package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"murim-helper/internal/domain"
	"murim-helper/internal/repository"
	"murim-helper/internal/service"
	"strings"
	"time"
)

type ScheduleUsecase interface {
	GenerateSchedule(ctx context.Context, description string) ([]domain.Schedule, error)
	UpdateSchedule(ctx context.Context, id string, updated domain.Schedule) error
	GetAllSchedules(ctx context.Context) ([]domain.Schedule, error)
	GetScheduleByID(ctx context.Context, id string) (*domain.Schedule, error)
	GetTodaySchedules(ctx context.Context) ([]domain.Schedule, error)
	GetThisWeekSchedules(ctx context.Context) ([]domain.Schedule, error)
	DeleteScheduleByID(ctx context.Context, id string) error
	MarkScheduleAsDone(ctx context.Context, id string) error
	MarkScheduleAsUndone(ctx context.Context, id string) error
	DeleteAll(ctx context.Context) error
}

type scheduleUsecase struct {
	repo   *repository.PostgresRepo
	groqAI service.GroqService
}

func NewScheduleUsecase(r *repository.PostgresRepo, ai service.GroqService) ScheduleUsecase {
	return &scheduleUsecase{repo: r, groqAI: ai}
}

func (s *scheduleUsecase) GenerateSchedule(ctx context.Context, desc string) ([]domain.Schedule, error) {
	if strings.TrimSpace(desc) == "" {
		return nil, errors.New("description cannot be empty")
	}

	// Add timeout for AI call
	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	schedules, err := s.groqAI.GenerateScheduleFromText(desc)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schedule from text: %w", err)
	}

	if len(schedules) == 0 {
		return nil, errors.New("AI returned no schedules")
	}

	if err := s.repo.SaveMany(ctx, schedules); err != nil {
		return nil, fmt.Errorf("failed to save generated schedules: %w", err)
	}

	return schedules, nil
}

func (s *scheduleUsecase) UpdateSchedule(ctx context.Context, id string, updated domain.Schedule) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id cannot be empty")
	}

	// Check if exists
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("schedule with id %s not found", id)
		}
		return fmt.Errorf("failed to check schedule existence: %w", err)
	}

	// Merge fields (partial update support)
	if updated.Title == "" {
		updated.Title = existing.Title
	}
	if updated.Description == "" {
		updated.Description = existing.Description
	}
	if updated.StartTime.IsZero() {
		updated.StartTime = existing.StartTime
	}
	if updated.EndTime.IsZero() {
		updated.EndTime = existing.EndTime
	}
	if updated.RepeatType == "" {
		updated.RepeatType = existing.RepeatType
	}
	if updated.RepeatUntil.IsZero() && !existing.RepeatUntil.IsZero() {
		updated.RepeatUntil = existing.RepeatUntil
	}

	return s.repo.Update(ctx, id, updated)
}

func (s *scheduleUsecase) GetAllSchedules(ctx context.Context) ([]domain.Schedule, error) {
	return s.repo.GetAll(ctx)
}

func (s *scheduleUsecase) GetScheduleByID(ctx context.Context, id string) (*domain.Schedule, error) {
	if strings.TrimSpace(id) == "" {
		return nil, errors.New("id cannot be empty")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *scheduleUsecase) GetTodaySchedules(ctx context.Context) ([]domain.Schedule, error) {
	loc := time.FixedZone("Asia/Jakarta", 7*3600)
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	return s.repo.GetTodaySchedules(ctx, startOfDay.Format(time.RFC3339), endOfDay.Format(time.RFC3339))
}

func (s *scheduleUsecase) GetThisWeekSchedules(ctx context.Context) ([]domain.Schedule, error) {
	loc := time.FixedZone("Asia/Jakarta", 7*3600)
	now := time.Now().In(loc)

	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -weekday+1)
	startOfWeek := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, loc)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return s.repo.GetThisWeekSchedules(ctx, startOfWeek.Format(time.RFC3339), endOfWeek.Format(time.RFC3339))
}

func (s *scheduleUsecase) DeleteScheduleByID(ctx context.Context, id string) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id cannot be empty")
	}

	// Check existence
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("schedule with id %s not found", id)
		}
		return fmt.Errorf("failed to check schedule existence: %w", err)
	}

	return s.repo.DeleteByID(ctx, id)
}

func (s *scheduleUsecase) MarkScheduleAsDone(ctx context.Context, id string) error {
	return s.setDoneStatus(ctx, id, true)
}

func (s *scheduleUsecase) MarkScheduleAsUndone(ctx context.Context, id string) error {
	return s.setDoneStatus(ctx, id, false)
}

func (s *scheduleUsecase) setDoneStatus(ctx context.Context, id string, done bool) error {
	if strings.TrimSpace(id) == "" {
		return errors.New("id cannot be empty")
	}

	schedule, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("schedule with id %s not found", id)
		}
		return fmt.Errorf("failed to get schedule by ID: %w", err)
	}

	schedule.IsDone = done
	return s.repo.Update(ctx, id, *schedule)
}

func (s *scheduleUsecase) DeleteAll(ctx context.Context) error {
	return s.repo.DeleteAll(ctx)
}
