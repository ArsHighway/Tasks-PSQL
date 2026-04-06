package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ArsHighway/Tasks-PSQL/internal/errs"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type taskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *taskRepository {
	return &taskRepository{pool: pool}
}

type TaskRepository interface {
	GetTaskWithID(ctx context.Context, id int) (*models.Task, error)
	GetTasks(ctx context.Context, args []any, baseQuery string) ([]models.Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error)
	CreateTask(ctx context.Context, t *models.Task) (*models.Task, error)
	UpdateTask(ctx context.Context, id int, t *models.Task) (*models.Task, error)
	PatchTask(ctx context.Context, id int, updates map[string]interface{}, parts []string, arg []interface{}) (*models.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

func (r *taskRepository) CreateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	err := r.pool.QueryRow(
		ctx,
		`INSERT INTO tasks (title, description, status, user_id, created_at)
		 VALUES ($1,$2,$3,$4,$5) RETURNING id`,
		t.Title,
		t.Description,
		t.Status,
		t.UserID,
		t.CreatedAt,
	).Scan(&t.ID)

	if err != nil {
		return nil, err
	}
	return t, nil
}

func (r *taskRepository) GetTaskWithID(ctx context.Context, id int) (*models.Task, error) {
	var t models.Task
	err := r.pool.QueryRow(ctx, `SELECT title,description,status FROM tasks WHERE id = $1`, id).Scan(
		&t.Title, &t.Description, &t.Status,
	)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, id int, t *models.Task) (*models.Task, error) {
	cmdTag, err := r.pool.Exec(ctx, `UPDATE tasks SET title =$1,description =$2,status =$3 WHERE id=$4`, t.Title, t.Description, t.Status, id)
	if err != nil {
		return nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, errs.ErrTaskNotFound
	}
	if err := r.pool.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1`, id).Scan(&t.ID, &t.Title,
		&t.Status, &t.UserID, &t.CreatedAt); err != nil {
		return nil, errs.ErrInvalidTask
	}
	return t, nil
}

func (r *taskRepository) PatchTask(ctx context.Context, id int, updates map[string]interface{}, parts []string, arg []interface{}) (*models.Task, error) {
	var t models.Task
	if len(parts) == 0 {
		return nil, errs.ErrNotValidFields
	}
	sql := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d", strings.Join(parts, ", "), id)
	cmdTag, err := r.pool.Exec(ctx, sql, arg...)
	if err != nil {
		return nil, errs.ErrInvalidTask
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, errs.ErrTaskNotFound
	}
	if err := r.pool.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1`, id).Scan(&t.ID, &t.Title,
		&t.Status, &t.UserID, &t.CreatedAt); err != nil {
		return nil, errs.ErrInvalidTask
	}
	return &t, nil
}

func (r *taskRepository) GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, title, description, status, user_id, created_at FROM tasks WHERE user_id = $1 ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.UserID, &t.CreatedAt); err != nil {
			return nil, errs.ErrInvalidTask
		}
		tasks = append(tasks, t)
	}
	return tasks, nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errs.ErrTaskNotFound
	}
	return nil
}

func (r *taskRepository) GetTasks(ctx context.Context, args []any, baseQuery string) ([]models.Task, error) {
	rows, err := r.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		return nil, errs.ErrTaskNotFound
	}
	defer rows.Close()
	tasks := []models.Task{}
	for rows.Next() {
		var t models.Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.UserID,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, errs.ErrInvalidTask
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		return nil, errs.ErrNotValidFields
	}
	return tasks, nil
}
