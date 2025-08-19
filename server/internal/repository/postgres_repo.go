package repository

import (
	"context"
	"database/sql"
	"fmt"
	"murim-helper/internal/domain"
	"murim-helper/internal/dto"
	"strings"
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

func (r *PostgresRepo) GetAll(ctx context.Context, page, limit int, filter dto.ScheduleFilter) ([]domain.Schedule, int, error) {
	var schedules []domain.Schedule
	var args []interface{}
	var conditions []string

	if filter.IsDone != nil {
		args = append(args, *filter.IsDone)
		conditions = append(conditions, fmt.Sprintf("is_done = $%d", len(args)))
	}
	if filter.RepeatType != "" {
		args = append(args, filter.RepeatType)
		conditions = append(conditions, fmt.Sprintf("repeat_type = $%d", len(args)))
	}
	if filter.Search != "" {
		args = append(args, "%"+filter.Search+"%")
		args = append(args, "%"+filter.Search+"%")
		conditions = append(conditions, fmt.Sprintf("(title ILIKE $%d OR description ILIKE $%d)", len(args)-1, len(args)))
	}
	if filter.StartAfter != nil {
		args = append(args, *filter.StartAfter)
		conditions = append(conditions, fmt.Sprintf("start_time >= $%d", len(args)))
	}
	if filter.StartBefore != nil {
		args = append(args, *filter.StartBefore)
		conditions = append(conditions, fmt.Sprintf("start_time < $%d", len(args)))
	}

	query := `SELECT * FROM schedules`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Sorting whitelist
	allowedSortColumns := map[string]bool{
		"start_time": true,
		"end_time":   true,
		"created_at": true,
		"title":      true,
	}
	sortBy := "start_time" // default
	if allowedSortColumns[filter.SortBy] {
		sortBy = filter.SortBy
	}
	sortOrder := "ASC"
	if strings.ToLower(filter.SortOrder) == "desc" {
		sortOrder = "DESC"
	}
	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Pagination
	offset := (page - 1) * limit
	args = append(args, limit, offset)
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", len(args)-1, len(args))

	if err := r.db.SelectContext(ctx, &schedules, query, args...); err != nil {
		return nil, 0, fmt.Errorf("failed to fetch schedules: %w", err)
	}

	// Count query
	countQuery := `SELECT COUNT(*) FROM schedules`
	if len(conditions) > 0 {
		countQuery += " WHERE " + strings.Join(conditions, " AND ")
	}
	var total int
	if err := r.db.GetContext(ctx, &total, countQuery, args[:len(args)-2]...); err != nil {
		return nil, 0, fmt.Errorf("failed to count schedules: %w", err)
	}

	return schedules, total, nil
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

func (r *PostgresRepo) GetRepeatingSchedules(ctx context.Context) ([]domain.Schedule, error) {
	var schedules []domain.Schedule
	query := `
        SELECT * FROM schedules
        WHERE repeat_type != 'none'
        AND (repeat_until IS NULL OR repeat_until > NOW())
        ORDER BY start_time ASC
    `
	if err := r.db.SelectContext(ctx, &schedules, query); err != nil {
		return nil, fmt.Errorf("failed to fetch repeating schedules: %w", err)
	}
	return schedules, nil
}

func (r *PostgresRepo) ExistsByStartTime(ctx context.Context, title string, start time.Time) (bool, error) {
	var count int
	err := r.db.GetContext(ctx, &count, `
        SELECT COUNT(*) FROM schedules
        WHERE title = $1 AND start_time = $2
    `, title, start)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
