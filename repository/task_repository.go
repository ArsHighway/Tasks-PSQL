package repository

import (
	"context"
	"encoding/json"
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
	ID          int
	Title       string
	Description string
	Status      string
	UserID      int
	CreatedAt   time.Time
}

type TaskRepository struct {
	pool *pgxpool.Pool
}

func NewTaskRepository(pool *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{pool: pool}
}

func (r *TaskRepository) CreateTask(ctx context.Context, w http.ResponseWriter, t *Task, log slog.Logger) error {
	log.Info("Creating task", "task", t.Title, "status", t.Status)
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	var id int
	err := r.pool.QueryRow(ctx, `INSERT INTO tasks (title,description,status,user_id,created_at) VALUES ($1,$2,$3,$4,$5) RETURNING id`, t.Title, t.Description, t.Status, t.UserID, t.CreatedAt).Scan(&id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Warn("DB insert failed", "error", err)
		return err
	}
	t.ID = id
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return err
	}
	log.Info("Task created successfully", "id", t.ID)
	return nil
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
