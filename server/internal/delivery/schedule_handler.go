package delivery

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"murim-helper/internal/dto"
	"murim-helper/internal/usecase"
	"murim-helper/pkg/httphelper"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type ScheduleHandler struct {
	Usecase usecase.ScheduleUsecase
}

func NewScheduleHandler(r *mux.Router, uc usecase.ScheduleUsecase) {
	handler := &ScheduleHandler{Usecase: uc}

	s := r.PathPrefix("/schedule").Subrouter()
	s.HandleFunc("", handler.Generate).Methods("POST")
	s.HandleFunc("", handler.GetAll).Methods("GET")
	s.HandleFunc("", handler.DeleteAll).Methods("DELETE")

	s.HandleFunc("/today", handler.GetToday).Methods("GET")
	s.HandleFunc("/this-week", handler.GetThisWeek).Methods("GET")

	s.HandleFunc("/{id}", handler.Update).Methods("PUT")
	s.HandleFunc("/{id}", handler.GetByID).Methods("GET")
	s.HandleFunc("/{id}", handler.DeleteByID).Methods("DELETE")
	s.HandleFunc("/{id}/done", handler.MarkAsDone).Methods("PUT")
	s.HandleFunc("/{id}/undone", handler.MarkAsUndone).Methods("PUT")
}

// ======================
// Helpers
// ======================

func withTimeout(r *http.Request, d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), d)
}

func getIDParam(r *http.Request) string {
	return mux.Vars(r)["id"]
}

func parsePagination(r *http.Request) (page, limit int) {
	page, _ = strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 10
	}
	return
}

// ======================
// Handlers
// ======================

func (h *ScheduleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	// Pagination support
	// page, limit := parsePagination(r)

	result, err := h.Usecase.GetAllSchedules(ctx /* optionally pass page, limit */)
	if err != nil {
		log.Printf("[GetAll] error: %v", err)
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch schedules", 50001)
		return
	}
	response := dto.ToScheduleResponseDTOs(result)
	httphelper.Success(w, r, http.StatusOK, "Successfully fetched schedules", response)
}

func (h *ScheduleHandler) Generate(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 15*time.Second) // longer for AI
	defer cancel()

	var req dto.GenerateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelper.Error(w, r, http.StatusBadRequest, "Invalid request body", 40001)
		return
	}

	if err := req.Validate(); err != nil {
		httphelper.Error(w, r, http.StatusBadRequest, err.Error(), 40002)
		return
	}

	result, err := h.Usecase.GenerateSchedule(ctx, req.Description)
	if err != nil {
		log.Printf("[Generate] error: %v", err)
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to generate schedule", 50002)
		return
	}
	httphelper.Success(w, r, http.StatusCreated, "Successfully generated schedule", dto.ToScheduleResponseDTOs(result))
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	id := getIDParam(r)
	if strings.TrimSpace(id) == "" {
		httphelper.Error(w, r, http.StatusBadRequest, "ID is required", 40003)
		return
	}

	var req dto.UpdateScheduleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		httphelper.Error(w, r, http.StatusBadRequest, "Invalid request body", 40004)
		return
	}

	if err := req.Validate(); err != nil {
		httphelper.Error(w, r, http.StatusBadRequest, err.Error(), 40005)
		return
	}

	// Get existing schedule for merging
	existing, err := h.Usecase.GetScheduleByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40402)
			return
		}
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch schedule", 50003)
		return
	}

	updated := req.ToDomain(*existing)
	if err := h.Usecase.UpdateSchedule(ctx, id, updated); err != nil {
		log.Printf("[Update] error: %v", err)
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to update schedule", 50004)
		return
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully updated schedule", nil)
}

func (h *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	id := getIDParam(r)
	if strings.TrimSpace(id) == "" {
		httphelper.Error(w, r, http.StatusBadRequest, "ID is required", 40006)
		return
	}

	result, err := h.Usecase.GetScheduleByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40401)
			return
		}
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch schedule", 50005)
		return
	}
	response := dto.ToScheduleResponseDTO(*result)
	httphelper.Success(w, r, http.StatusOK, "Successfully fetched schedule", response)
}

func (h *ScheduleHandler) GetToday(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	schedules, err := h.Usecase.GetTodaySchedules(ctx)
	if err != nil {
		log.Printf("[GetToday] error: %v", err)
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch today's schedules", 50006)
		return
	}

	httphelper.Success(w, r, http.StatusOK, "Successfully fetched today's schedule", dto.ToScheduleResponseDTOs(schedules))
}

func (h *ScheduleHandler) GetThisWeek(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	schedules, err := h.Usecase.GetThisWeekSchedules(ctx)
	if err != nil {
		log.Printf("[GetThisWeek] error: %v", err)
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to fetch this week's schedules", 50007)
		return
	}

	httphelper.Success(w, r, http.StatusOK, "Successfully fetched this week's schedules", dto.ToScheduleResponseDTOs(schedules))
}

func (h *ScheduleHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	id := getIDParam(r)
	if strings.TrimSpace(id) == "" {
		httphelper.Error(w, r, http.StatusBadRequest, "ID is required", 40007)
		return
	}

	err := h.Usecase.DeleteScheduleByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40403)
			return
		}
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to delete schedule", 50008)
		return
	}
	httphelper.Success(w, r, http.StatusNoContent, "Successfully deleted schedule", nil)
}

func (h *ScheduleHandler) MarkAsDone(w http.ResponseWriter, r *http.Request) {
	h.setDoneStatus(w, r, true)
}

func (h *ScheduleHandler) MarkAsUndone(w http.ResponseWriter, r *http.Request) {
	h.setDoneStatus(w, r, false)
}

func (h *ScheduleHandler) setDoneStatus(w http.ResponseWriter, r *http.Request, done bool) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	id := getIDParam(r)
	if strings.TrimSpace(id) == "" {
		httphelper.Error(w, r, http.StatusBadRequest, "ID is required", 40008)
		return
	}

	var err error
	if done {
		err = h.Usecase.MarkScheduleAsDone(ctx, id)
	} else {
		err = h.Usecase.MarkScheduleAsUndone(ctx, id)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			httphelper.Error(w, r, http.StatusNotFound, "Schedule not found", 40404)
			return
		}
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to update schedule status", 50009)
		return
	}

	statusMsg := "undone"
	if done {
		statusMsg = "done"
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully marked schedule as "+statusMsg, nil)
}

func (h *ScheduleHandler) DeleteAll(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := withTimeout(r, 5*time.Second)
	defer cancel()

	err := h.Usecase.DeleteAll(ctx)
	if err != nil {
		log.Printf("[DeleteAll] error: %v", err)
		httphelper.Error(w, r, http.StatusInternalServerError, "Failed to delete all schedules", 50010)
		return
	}
	httphelper.Success(w, r, http.StatusOK, "Successfully deleted all schedules", nil)
}
