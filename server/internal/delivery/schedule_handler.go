package delivery

import (
	"encoding/json"
	"fmt"
	"net/http"

	"murim-helper/internal/domain"
	"murim-helper/internal/dto"
	"murim-helper/internal/usecase"
	"murim-helper/pkg/httphelper"

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
	r.HandleFunc("/schedule/today", handler.GetToday).Methods("GET")
	r.HandleFunc("/schedule/this-week", handler.GetThisWeek).Methods("GET")
	r.HandleFunc("/schedule/{id}", handler.GetByID).Methods("GET")
	r.HandleFunc("/schedule/{id}", handler.DeleteByID).Methods("DELETE")
	r.HandleFunc("/schedule/{id}/done", handler.MarkAsDone).Methods("PUT")
	r.HandleFunc("/schedule/{id}/undone", handler.MarkAsUndone).Methods("PUT")
	r.HandleFunc("/schedule", handler.DeleteAll).Methods("DELETE")
}

func (h *ScheduleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	result, err := h.Usecase.GetAllSchedules()
	if err != nil {
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch schedules", 50001)
		return
	}
	response := dto.ToScheduleResponseDTOs(result)
	httphelper.Success(w, r, http.StatusOK, "Successfully fetched schedules", response)
}

func (h *ScheduleHandler) Generate(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		httphelper.Error(w, r, http.StatusBadRequest, "Invalid request body", 40001)
		return
	}

	result, err := h.Usecase.GenerateSchedule(body.Description)
	if err != nil {
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to generate schedule", 50002)
		return
	}
	httphelper.Success(w, r, http.StatusCreated, "Successfully generated schedule", result)
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	var updated domain.Schedule
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		httphelper.Error(w, r, http.StatusBadRequest, "Invalid request body", 40002)
		return
	}

	err := h.Usecase.UpdateSchedule(id, updated)
	if err != nil {
		httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40402)
		return
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully updated schedule", nil)
}

func (h *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result, err := h.Usecase.GetScheduleByID(id)
	if err != nil {
		httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40401)
		return
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully fetched schedule", result)
}

func (h *ScheduleHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.Usecase.GetTodaySchedules()
	if err != nil {
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch today's schedules", 50002)
		return
	}

	httphelper.Success(w, r, http.StatusOK, "Successfully fetched today's schedule", schedules)
}

func (h *ScheduleHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	err := h.Usecase.DeleteScheduleByID(id)
	if err != nil {
		httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40403)
		return
	}
	httphelper.Success(w, r, http.StatusNoContent, "Successfully deleted schedule", nil)
}

func (h *ScheduleHandler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Marking schedule as done")
	id := mux.Vars(r)["id"]
	err := h.Usecase.MarkScheduleAsDone(id)
	if err != nil {
		httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40404)
		return
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully marked schedule as done", nil)
}

func (h *ScheduleHandler) MarkAsUndone(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Marking schedule as undone")
	id := mux.Vars(r)["id"]
	err := h.Usecase.MarkScheduleAsUndone(id) // Reuse the same method for toggling done state
	if err != nil {
		httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40405)
		return
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully marked schedule as undone", nil)
}

func (h *ScheduleHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	err := h.Usecase.DeleteAll()
	if err != nil {
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to delete all schedules", 50003)
		return
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully deleted all schedules", nil)
}

func (h *ScheduleHandler) GetThisWeek(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.Usecase.GetThisWeekSchedules()
	if err != nil {
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch this week's schedules", 50004)
		return
	}

	httphelper.Success(w, r, http.StatusOK, "Successfully fetched this week's schedules", schedules)
}
