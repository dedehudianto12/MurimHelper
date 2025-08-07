package usecase

import (
	"fmt"

	"murim-helper/internal/domain"
	"murim-helper/internal/repository"
	"murim-helper/internal/service"
)

type ScheduleUsecase interface {
	GenerateSchedule(description string) ([]domain.Schedule, error)
	UpdateSchedule(id string, updated domain.Schedule) error
	GetAllSchedules() ([]domain.Schedule, error)
	GetScheduleByID(id string) (*domain.Schedule, error)
	GetTodaySchedules() ([]domain.Schedule, error)
	GetThisWeekSchedules() ([]domain.Schedule, error)
	DeleteScheduleByID(id string) error
	MarkScheduleAsDone(id string) error
	MarkScheduleAsUndone(id string) error
	DeleteAll() error
}

type scheduleUsecase struct {
	repo *repository.SQLiteRepo
	// openai service.OpenAIService
	// ollamaAI service.OllamaService
	groqAI service.GroqService
}

func NewScheduleUsecase(r *repository.SQLiteRepo, ai service.GroqService) ScheduleUsecase {
	return &scheduleUsecase{repo: r, groqAI: ai}
}

func (s *scheduleUsecase) GenerateSchedule(desc string) ([]domain.Schedule, error) {
	schedules, err := s.groqAI.GenerateScheduleFromText(desc)
	if err != nil {
		return nil, fmt.Errorf("failed to generate schedule from text: %w", err)
	}

	err = s.repo.SaveMany(schedules)
	if err != nil {
		return nil, fmt.Errorf("failed to save generated schedules: %w", err)
	}

	return schedules, nil
}

func (s *scheduleUsecase) UpdateSchedule(id string, updated domain.Schedule) error {
	return s.repo.Update(id, updated)
}

func (s *scheduleUsecase) GetAllSchedules() ([]domain.Schedule, error) {
	return s.repo.GetAll()
}

func (s *scheduleUsecase) GetScheduleByID(id string) (*domain.Schedule, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}
	return schedule, nil
}

func (s *scheduleUsecase) GetTodaySchedules() ([]domain.Schedule, error) {
	return s.repo.GetTodaySchedules()
}

func (s *scheduleUsecase) GetThisWeekSchedules() ([]domain.Schedule, error) {
	return s.repo.GetThisWeekSchedules()
}

func (s *scheduleUsecase) DeleteScheduleByID(id string) error {
	return s.repo.DeleteByID(id)
}

func (s *scheduleUsecase) MarkScheduleAsDone(id string) error {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get schedule by ID: %w", err)
	}
	schedule.IsDone = true
	return s.repo.Update(id, *schedule)
}

func (s *scheduleUsecase) MarkScheduleAsUndone(id string) error {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get schedule by ID: %w", err)
	}
	schedule.IsDone = false
	return s.repo.Update(id, *schedule)
}

func (s *scheduleUsecase) DeleteAll() error {
	return s.repo.DeleteAll()
}
