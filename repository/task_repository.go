package repository

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) CreateTask(ctx context.Context, t *Task, log slog.Logger) (*Task, error) {
	log.Info("Creating task",
		"title", t.Title,
		"user_id", t.UserID,
	)

	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}

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
		log.Error("DB insert failed", "error", err)
		return nil, err
	}

	log.Info("Task created", "taskID", t.ID)

	return t, nil
}

func (r *TaskRepository) GetTaskWithID(ctx context.Context, w http.ResponseWriter, id int, log slog.Logger) (*Task, error) {
	log.Info("Get task", "taskID", id)
	var t Task
	err := r.pool.QueryRow(ctx, `SELECT title,description,status FROM tasks WHERE id = $1`, id).Scan(
		&t.Title, &t.Description, &t.Status,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Warn("Task not found", "taskID", id)
			return nil, err
		}
		log.Error("DB query failed", "error", err)
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) UpdateTask(ctx context.Context, w http.ResponseWriter, id int, t *Task, log slog.Logger) (*Task, error) {
	log.Info("Update task", "taskID", id)
	cmdTag, err := r.pool.Exec(ctx, `UPDATE tasks SET title =$1,description =$2,status =$3 WHERE id=$4`, t.Title, t.Description, t.Status, id)
	if err != nil {
		log.Error("DB query failed", "error", err)
		return nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		log.Warn("Task not found", "taskID", id)
		return nil, pgx.ErrNoRows
	}
	if err := r.pool.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1`, id).Scan(&t.ID, &t.Title,
		&t.Status, &t.UserID, &t.CreatedAt); err != nil {
		log.Error("DB query after update failed", "error", err)
		return nil, err
	}
	return t, nil
}

func (r *TaskRepository) PatchTask(ctx context.Context, w http.ResponseWriter, id int, updates map[string]interface{}, log slog.Logger) (*Task, error) {
	log.Info("Patch task", "taskID", id)
	var arg []interface{}
	c := 1
	var t Task
	allowed := map[string]bool{
		"title":       true,
		"description": true,
		"status":      true,
	}
	parts := []string{}
	for k, v := range updates {
		if !allowed[k] {
			continue
		}
		arg = append(arg, v)
		parts = append(parts, fmt.Sprintf("%v =$%d", k, c))
		c++
	}
	if len(parts) == 0 {
		return nil, errors.New("no valid fields to update")
	}
	sql := fmt.Sprintf("UPDATE tasks SET %s WHERE id=$%d", strings.Join(parts, ", "), id)
	cmdTag, err := r.pool.Exec(ctx, sql, arg...)
	if err != nil {
		log.Error("DB query failed", "error", err)
		return nil, err
	}
	if cmdTag.RowsAffected() == 0 {
		log.Warn("Task not found", "taskID", id)
		return nil, pgx.ErrNoRows
	}
	if err := r.pool.QueryRow(ctx, `SELECT * FROM tasks WHERE id = $1`, id).Scan(&t.ID, &t.Title,
		&t.Status, &t.UserID, &t.CreatedAt); err != nil {
		log.Error("DB query after update failed", "error", err)
		return nil, err
	}
	return &t, nil
}

func (r *TaskRepository) DeleteTask(ctx context.Context, w http.ResponseWriter, id int, log slog.Logger) error {
	log.Info("Update task", "taskID", id)
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
	if err != nil {
		log.Error("DB query failed", "error", err)
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		log.Warn("Task not found", "taskID", id)
		return pgx.ErrNoRows
	}
	return nil
}

// func (r *TaskRepository) GetTasks(ctx context.Context, status string, log slog.Logger) (*[]Task, error) {
// 	log.Info("Get tasks", "status", status)
// 	rows, err := r.pool.Query(ctx, `SELECT id,title,description,status,user_id FROM tasks WHERE status = $1`, status)
// 	if err != nil {
// 		log.Error("DB query failed", "error", err)
// 		return nil, err
// 	}
// 	var t Task
// 	var tasks []Task
// 	defer rows.Close()
// 	for rows.Next() {
// 		err := rows.Scan(
// 			&t.ID,
// 			&t.Title,
// 			&t.Description,
// 			&t.Status,
// 			&t.UserID,
// 		)
// 		tasks = append(tasks, t)
// 		if err != nil {
// 			log.Warn("Problem with scan", "error", err)
// 			return nil, err
// 		}
// 	}
// 	if len(tasks) == 0 {
// 		log.Warn("Task no found", "status", status)
// 		return nil, pgx.ErrNoRows
// 	}
// 	return &tasks, nil
// }

func (r *TaskRepository) GetTasks(ctx context.Context, args []any, baseQuery string, log *slog.Logger) ([]Task, error) {
	log.Info("Get tasks")
	rows, err := r.pool.Query(ctx, baseQuery, args...)
	if err != nil {
		log.Error("DB query failed", "error", err)
		return nil, err
	}
	defer rows.Close()
	tasks := []Task{}
	for rows.Next() {
		var t Task
		err := rows.Scan(
			&t.ID,
			&t.Title,
			&t.Description,
			&t.Status,
			&t.UserID,
			&t.CreatedAt,
		)
		if err != nil {
			log.Warn("Failed to scan tasks")
			return nil, err
		}
		tasks = append(tasks, t)
	}
	if len(tasks) == 0 {
		log.Warn("Tasks no found")
		return nil, pgx.ErrNoRows
	}
	return tasks, nil
}
