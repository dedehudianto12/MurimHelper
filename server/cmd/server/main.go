package main

import (
	"log"
	"net/http"

	"murim-helper/internal/delivery"
	"murim-helper/internal/repository"
	"murim-helper/internal/service"
	"murim-helper/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	repo, err := repository.NewSQLiteRepo("schedules.db")
	if err != nil {
		log.Fatalf("Failed connect to DB: %v", err)
	}
	ai := service.NewOllamaService()
	uc := usecase.NewScheduleUsecase(repo, ai)

	r := mux.NewRouter()
	delivery.NewScheduleHandler(r, uc)

	if err != nil {
		log.Fatalf("Failed to create repository: %v", err)
	}
	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
