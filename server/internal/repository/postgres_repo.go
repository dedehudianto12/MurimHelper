package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murim-helper/internal/domain"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PostgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(connStr string) (*PostgresRepo, error) {
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}
	return &PostgresRepo{db: db}, nil
}

// SaveMany inserts multiple schedules in one batch inside a transaction
func (r *PostgresRepo) SaveMany(ctx context.Context, schedules []domain.Schedule) error {
	if len(schedules) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction failed: %w", err)
	}

	query := `INSERT INTO schedules 
		(id, title, description, start_time, end_time, is_done, repeat_type, repeat_until) VALUES `

	args := []interface{}{}
	placeholders := []string{}

	for i, s := range schedules {
		idx := i * 8
		placeholders = append(placeholders,
			fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)",
				idx+1, idx+2, idx+3, idx+4, idx+5, idx+6, idx+7, idx+8))
		args = append(args,
			s.ID, s.Title, s.Description, s.StartTime, s.EndTime, s.IsDone, s.RepeatType, s.RepeatUntil)
	}

	query += strings.Join(placeholders, ",")
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		tx.Rollback()
		return fmt.Errorf("batch insert failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}
	return nil
}

func (r *PostgresRepo) Update(ctx context.Context, id string, updated domain.Schedule) error {
	query := `
		UPDATE schedules
		SET title = $1, description = $2, start_time = $3, end_time = $4, is_done = $5, repeat_type = $6, repeat_until = $7
		WHERE id = $8`

	res, err := r.db.ExecContext(ctx, query,
		updated.Title, updated.Description, updated.StartTime, updated.EndTime,
		updated.IsDone, updated.RepeatType, updated.RepeatUntil, id)
	if err != nil {
		return fmt.Errorf("update failed: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresRepo) GetAll(ctx context.Context) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	err := r.db.SelectContext(ctx, &schedules, `SELECT * FROM schedules ORDER BY start_time ASC`)
	if err != nil {
		return nil, fmt.Errorf("get all failed: %w", err)
	}
	return schedules, nil
}

func (r *PostgresRepo) GetByID(ctx context.Context, id string) (*domain.Schedule, error) {
	var schedule domain.Schedule
	err := r.db.GetContext(ctx, &schedule, `SELECT * FROM schedules WHERE id = $1`, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, fmt.Errorf("get by id failed: %w", err)
	}
	return &schedule, nil
}

func (r *PostgresRepo) GetTodaySchedules(ctx context.Context, startOfDay, endOfDay string) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	query := `
		SELECT * FROM schedules
		WHERE start_time >= $1 AND start_time < $2
		ORDER BY start_time ASC`
	err := r.db.SelectContext(ctx, &schedules, query, startOfDay, endOfDay)
	if err != nil {
		return nil, fmt.Errorf("get today schedules failed: %w", err)
	}
	return schedules, nil
}

func (r *PostgresRepo) GetThisWeekSchedules(ctx context.Context, startOfWeek, endOfWeek string) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	query := `
		SELECT * FROM schedules
		WHERE start_time >= $1 AND start_time < $2
		ORDER BY start_time ASC`
	err := r.db.SelectContext(ctx, &schedules, query, startOfWeek, endOfWeek)
	if err != nil {
		return nil, fmt.Errorf("get this week schedules failed: %w", err)
	}
	return schedules, nil
}

func (r *PostgresRepo) DeleteByID(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM schedules WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("delete by id failed: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (r *PostgresRepo) DeleteAll(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM schedules`)
	if err != nil {
		return fmt.Errorf("delete all failed: %w", err)
	}
	return nil
}
