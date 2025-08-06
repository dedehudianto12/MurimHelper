package repository

import (
	"errors"
	"murim-helper/internal/model"
)

type InMemoryRepo struct {
	data map[string]model.Schedule
}

func NewInMemoryRepo() *InMemoryRepo {
	return &InMemoryRepo{data: make(map[string]model.Schedule)}
}

func (r *InMemoryRepo) SaveMany(schedules []model.Schedule){
	for _, s := range schedules {
		r.data[s.ID] = s
	}
}

func (r *InMemoryRepo) Update(id string, updated model.Schedule) error {
    if _, exists := r.data[id]; !exists {
        return errors.New("not found")
    }
    r.data[id] = updated
    return nil
}

func (r *InMemoryRepo) GetAll() []model.Schedule{ 
	var result []model.Schedule
	for _, v := range r.data {
		result = append(result, v)
	}
	return result
}