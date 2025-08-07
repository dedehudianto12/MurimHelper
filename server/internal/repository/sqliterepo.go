package repository

import (
	"fmt"
	"murim-helper/internal/domain"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRepo struct {
	db *sqlx.DB
}

func NewSQLiteRepo(dbPath string) (*SQLiteRepo, error) {
	db, err := sqlx.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS schedules (
		id TEXT PRIMARY KEY,
		title TEXT NOT NULL,
		description TEXT,
		start_time DATETIME NOT NULL,
		end_time DATETIME NOT NULL,
		is_done BOOLEAN DEFAULT FALSE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	db.MustExec(schema)

	return &SQLiteRepo{db: db}, nil
}

func (r *SQLiteRepo) SaveMany(schedules []domain.Schedule) error {
	query := `INSERT INTO schedules 
	(id, title, description, start_time, end_time, is_done) 
	VALUES (?, ?, ?, ?, ?, ?)`

	for _, s := range schedules {
		_, err := r.db.Exec(query,
			s.ID, s.Title, s.Description, s.StartTime, s.EndTime, s.IsDone)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *SQLiteRepo) Update(id string, updated domain.Schedule) error {
	query := `
	UPDATE schedules
	SET title = ?, description = ?, start_time = ?, end_time = ?, is_done = ?
	WHERE id = ?`

	_, err := r.db.Exec(query,
		updated.Title, updated.Description, updated.StartTime, updated.EndTime, updated.IsDone, id)
	return err
}

func (r *SQLiteRepo) GetAll() ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.Select(&schedules, `SELECT * FROM schedules ORDER BY start_time ASC`)
	return schedules, err
}

func (r *SQLiteRepo) GetByID(id string) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := r.db.Get(&schedule, `SELECT * FROM schedules WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *SQLiteRepo) GetTodaySchedules() ([]domain.Schedule, error) {
	loc := time.FixedZone("Asia/Jakarta", 7*3600)
	now := time.Now().In(loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	endOfDay := startOfDay.Add(24 * time.Hour)

	query := `
	SELECT id, title, description, start_time, end_time, is_done, created_at
	FROM schedules
	WHERE start_time >= ? AND start_time < ?
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

func (r *SQLiteRepo) GetThisWeekSchedules() ([]domain.Schedule, error) {
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
	WHERE start_time >= ? AND start_time < ?
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

func (r *SQLiteRepo) DeleteByID(id string) error {
	_, err := r.db.Exec(`DELETE FROM schedules WHERE id = ?`, id)
	return err
}

func (r *SQLiteRepo) DeleteAll() error {
	_, err := r.db.Exec(`DELETE FROM schedules`)
	return err
}
