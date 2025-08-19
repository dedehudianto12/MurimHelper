// @title Murim Helper API
// @version 1.0
// @description This is the API documentation for Murim Helper (schedule/todo app).
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@murimhelper.local

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
package main

import (
	"log"
	"net/http"
	"os"

	"murim-helper/internal/delivery"
	"murim-helper/internal/repository"
	"murim-helper/internal/service"
	"murim-helper/internal/service/cronjob"
	"murim-helper/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "murim-helper/docs" // 👈 import generated docs

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
	cronjob.StartCronJobs(uc)

	r := mux.NewRouter()
	delivery.NewScheduleHandler(r, uc)

	r.PathPrefix("/docs/").Handler(httpSwagger.WrapHandler)

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
