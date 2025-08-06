package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"murim-helper/internal/model"
	"murim-helper/internal/usecase"

	"github.com/gorilla/mux"
)

type ScheduleHandler struct {
	Usecase usecase.ScheduleUsecase
}

func NewScheduleHandler(r *mux.Router, uc usecase.ScheduleUsecase) {
	handler := &ScheduleHandler{Usecase: uc}
	r.HandleFunc("/schedule", handler.Generate).Methods("POST")
	r.HandleFunc("/schedule/{id}", handler.Update).Methods("PUT")
	r.HandleFunc("/schedule", handler.GetAll).Methods("GET")
	r.HandleFunc("/schedule/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/schedule/{id}", handler.DeleteByID).Methods("DELETE")
	r.HandleFunc("/schedule/{id}/done", handler.MarkAsDone).Methods("PUT")
}

func (h *ScheduleHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Description string `json:"description"`
	}
	json.NewDecoder(r.Body).Decode(&body)

	result, err := h.Usecase.GenerateSchedule(body.Description)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updated model.Schedule
	json.NewDecoder(r.Body).Decode(&updated)

	err := h.Usecase.UpdateSchedule(id, updated)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *ScheduleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, _ := h.Usecase.GetAllSchedules()
	json.NewEncoder(w).Encode(result)
}

func (h *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result, err := h.Usecase.GetScheduleByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(result)
}

func (h *ScheduleHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := h.Usecase.DeleteScheduleByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ScheduleHandler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Marking schedule as done")
	id := mux.Vars(r)["id"]
	err := h.Usecase.MarkScheduleAsDone(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
