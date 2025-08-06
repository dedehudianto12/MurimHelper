package repository

import (
	"murim-helper/internal/model"

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

func (r *SQLiteRepo) SaveMany(schedules []model.Schedule) error {
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

func (r *SQLiteRepo) Update(id string, updated model.Schedule) error {
	query := `
	UPDATE schedules
	SET title = ?, description = ?, start_time = ?, end_time = ?, is_done = ?
	WHERE id = ?`

	_, err := r.db.Exec(query,
		updated.Title, updated.Description, updated.StartTime, updated.EndTime, updated.IsDone, id)
	return err
}

func (r *SQLiteRepo) GetAll() ([]model.Schedule, error) {
	var schedules []model.Schedule
	err := r.db.Select(&schedules, `SELECT * FROM schedules ORDER BY start_time ASC`)
	return schedules, err
}

func (r *SQLiteRepo) GetByID(id string) (*model.Schedule, error) {
	var schedule model.Schedule
	err := r.db.Get(&schedule, `SELECT * FROM schedules WHERE id = ?`, id)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

func (r *SQLiteRepo) DeleteByID(id string) error {
	_, err := r.db.Exec(`DELETE FROM schedules WHERE id = ?`, id)
	return err
}
