package main

import (
	"log"
	"net/http"
	"os"

	"murim-helper/internal/delivery"
	"murim-helper/internal/repository"
	"murim-helper/internal/service"
	"murim-helper/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	connStr := os.Getenv("POSTGRES_CONN") // e.g. "postgres://user:pass@localhost:5432/dbname?sslmode=disable"
	repo, err := repository.NewPostgresRepo(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	ai := service.NewGroqService()
	uc := usecase.NewScheduleUsecase(repo, ai)

	r := mux.NewRouter()
	delivery.NewScheduleHandler(r, uc)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
