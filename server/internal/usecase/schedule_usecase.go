package usecase

import (
	"errors"
	"murim-helper/internal/model"
	"murim-helper/internal/repository"

	"github.com/google/uuid"
)

type ScheduleUsecase interface {
    GenerateSchedule(desc string) ([]model.Schedule, error)
    UpdateSchedule(id string, updated model.Schedule) error
    GetAllSchedules() ([]model.Schedule, error)
	GetScheduleByID(id string) (*model.Schedule, error)
}

type scheduleUsecase struct {
	repo *repository.InMemoryRepo
}

func NewScheduleUsecase(r *repository.InMemoryRepo) ScheduleUsecase {
	return &scheduleUsecase{repo: r}
}

func (s *scheduleUsecase) GenerateSchedule(desc string) ([]model.Schedule, error) {
	// This is where OpenAI call would go later. For now, return dummy data
	schedules := []model.Schedule{
		{ID: uuid.NewString(), StartTime: "07:00", EndTime: "08:00",  Task: "Wake up and Read Bible", Description: desc},
		{ID: uuid.NewString(), StartTime: "08:00", EndTime: "09:00", Task: "Market Review", Description: desc},
	}
	s.repo.SaveMany(schedules)
	return schedules, nil
}

func (s *scheduleUsecase) UpdateSchedule(id string, updated model.Schedule) error {
	return s.repo.Update(id, updated)
}

func (s *scheduleUsecase) GetAllSchedules() ([]model.Schedule, error) {
	return s.repo.GetAll(), nil
}

func (s *scheduleUsecase) GetScheduleByID(id string) (*model.Schedule, error) {
	schedules := s.repo.GetAll()
	for _, schedule := range schedules {
		if schedule.ID == id {
			return &schedule, nil
		}
	}
	return nil, errors.New("schedule not found")
}