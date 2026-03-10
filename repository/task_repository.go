package repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

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
