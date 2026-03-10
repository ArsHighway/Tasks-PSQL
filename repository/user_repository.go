package repository

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) CreateUser(ctx context.Context, w http.ResponseWriter, u *User, log slog.Logger) error {
	log.Info("Creating user", "name", u.Name, "email", u.Email)
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	var id int
	err := r.pool.QueryRow(ctx, `INSERT INTO users (name,email,created_at) VALUES ($1, $2,$3) RETURNING id`, u.Name, u.Email, u.CreatedAt).Scan(&id)
	if err != nil {
		log.Warn("DB insert failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	u.ID = id
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding failed", "error", err)
		return err
	}
	log.Info("User created successfully", "id", u.ID)
	return nil
}

func (r *UserRepository) GetTaskWithUserID(ctx context.Context, w http.ResponseWriter, id int, log *slog.Logger) error {
	rows, err := r.pool.Query(ctx, `Select title,discription from tasks where user_id = $1`, id)
	if err != nil {
		log.Warn("DB insert failed", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return err
	}
	type t struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	var tasks []t
	defer rows.Close()
	for rows.Next() {
		var task t
		err := rows.Scan(
			&task.Title,
			&task.Description,
		)
		if err != nil {
			log.Warn("DB insert failed", "error", err)
			http.Error(w, "Problem with scan", http.StatusInternalServerError)
			return err
		}
		tasks = append(tasks, task)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return err
	}
	log.Info("user tasks found", "tasks", tasks)
	return nil
}
