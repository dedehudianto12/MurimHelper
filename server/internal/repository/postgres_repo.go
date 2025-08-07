package repository

import (
	"fmt"
	"murim-helper/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(connStr string) (*PostgresRepo, error) {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &PostgresRepo{db: db}, nil
}

func (r *PostgresRepo) SaveMany(schedules []domain.Schedule) error {
	query := `INSERT INTO schedules 
	(id, title, description, start_time, end_time, is_done) 
	VALUES ($1, $2, $3, $4, $5, $6)`

	for _, s := range schedules {
		_, err := r.db.Exec(query,
			s.ID, s.Title, s.Description, s.StartTime, s.EndTime, s.IsDone)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *PostgresRepo) Update(id string, updated domain.Schedule) error {
	query := `
	UPDATE schedules
	SET title = $1, description = $2, start_time = $3, end_time = $4, is_done = $5
	WHERE id = $6`

	_, err := r.db.Exec(query,
		updated.Title, updated.Description, updated.StartTime, updated.EndTime, updated.IsDone, id)
	return err
}

func (r *PostgresRepo) GetAll() ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.Select(&schedules, `SELECT * FROM schedules ORDER BY start_time ASC`)
	return schedules, err
}

func (r *PostgresRepo) GetByID(id string) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := r.db.Get(&schedule, `SELECT * FROM schedules WHERE id = $1`, id)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *PostgresRepo) GetTodaySchedules() ([]domain.Schedule, error) {
	loc := time.FixedZone("Asia/Jakarta", 7*3600)
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
	SELECT id, title, description, start_time, end_time, is_done, created_at
	FROM schedules
	WHERE start_time >= $1 AND start_time < $2
	ORDER BY start_time ASC
	`

	fmt.Println("Start of day:", startOfDay.Format(time.RFC3339))
	fmt.Println("End of day:", endOfDay.Format(time.RFC3339))

	rows, err := r.db.Query(query, startOfDay, endOfDay)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.Schedule
	for rows.Next() {
		var s domain.Schedule
		err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.StartTime, &s.EndTime, &s.IsDone, &s.CreatedAt)
		if err != nil {
			continue
		}
		schedules = append(schedules, s)
	}

	fmt.Printf("Fetched %d schedules for today\n", len(schedules))

	return schedules, nil
}

func (r *PostgresRepo) GetThisWeekSchedules() ([]domain.Schedule, error) {
	loc := time.FixedZone("Asia/Jakarta", 7*3600)
	now := time.Now().In(loc)

	// Find Monday of the current week
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday
		weekday = 7
	}
	monday := now.AddDate(0, 0, -weekday+1)
	startOfWeek := time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, loc)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	query := `
	SELECT id, title, description, start_time, end_time, is_done, created_at
	FROM schedules
	WHERE start_time >= $1 AND start_time < $2
	ORDER BY start_time ASC
	`

	rows, err := r.db.Query(query, startOfWeek, endOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []domain.Schedule
	for rows.Next() {
		var s domain.Schedule
		err := rows.Scan(&s.ID, &s.Title, &s.Description, &s.StartTime, &s.EndTime, &s.IsDone, &s.CreatedAt)
		if err != nil {
			continue
		}
		schedules = append(schedules, s)
	}

	return schedules, nil
}

func (r *PostgresRepo) DeleteByID(id string) error {
	_, err := r.db.Exec(`DELETE FROM schedules WHERE id = $1`, id)
	return err
}

func (r *PostgresRepo) DeleteAll() error {
	_, err := r.db.Exec(`DELETE FROM schedules`)
	return err
}
