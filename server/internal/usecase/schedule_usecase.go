package usecase

import (
	"fmt"

	"murim-helper/internal/model"
	"murim-helper/internal/repository"
	"murim-helper/internal/service"
)

type ScheduleUsecase interface {
	GenerateSchedule(description string) ([]model.Schedule, error)
	UpdateSchedule(id string, updated model.Schedule) error
	GetAllSchedules() ([]model.Schedule, error)
	GetScheduleByID(id string) (*model.Schedule, error)
	DeleteScheduleByID(id string) error
	MarkScheduleAsDone(id string) error
}

type scheduleUsecase struct {
	repo *repository.SQLiteRepo
	// openai service.OpenAIService
	// ollamaAI service.OllamaService
	groqAI service.GroqService
}

func NewScheduleUsecase(r *repository.SQLiteRepo, ai service.OllamaService) ScheduleUsecase {
	return &scheduleUsecase{repo: r, groqAI: ai}
}

func (s *scheduleUsecase) GenerateSchedule(desc string) ([]model.Schedule, error) {
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

func (s *scheduleUsecase) UpdateSchedule(id string, updated model.Schedule) error {
	return s.repo.Update(id, updated)
}

func (s *scheduleUsecase) GetAllSchedules() ([]model.Schedule, error) {
	return s.repo.GetAll()
}

func (s *scheduleUsecase) GetScheduleByID(id string) (*model.Schedule, error) {
	schedule, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get schedule by ID: %w", err)
	}
	return schedule, nil
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
