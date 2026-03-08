package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

type Tasks struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	UserID      int       `json:"user_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return UserRepository{pool: pool}
}

func (r *UserRepository) CreateUsers(ctx context.Context, w http.ResponseWriter, u *User, log slog.Logger) error {
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

func (r *UserRepository) TasksUser(ctx context.Context, w http.ResponseWriter, t *Tasks, log slog.Logger) error {
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

func (r *UserRepository) GetUserWithID(ctx context.Context, w http.ResponseWriter, id int, log *slog.Logger) error {
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

func HandlerCreateUser(w http.ResponseWriter, r *http.Request, ur *UserRepository) error {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	log := slog.With("handler", "CreateUsers",
		"request_method", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return errors.New("Bad method")
	}
	var u User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Problem with decoding", http.StatusNotFound)
		log.Warn("JSON decoding failed", "error", err)
		return err
	}
	return ur.CreateUsers(ctx, w, &u, *log)
}

func HandlerTasks(w http.ResponseWriter, r *http.Request, ur *UserRepository) error {
	ctx, cancle := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancle()
	log := slog.With("handler", "Tasks",
		"request_method", r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return errors.New("Bad method")
	}
	var t Tasks
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Problem with decoding", http.StatusNotFound)
		log.Warn("JSON decoding failed", "error", err)
		return err
	}
	return ur.TasksUser(ctx, w, &t, *log)
}

func HandlerGetUserWithID(w http.ResponseWriter, r *http.Request, ur *UserRepository) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	log := slog.With("handler", "GetUserWithID", "request_method", r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Сonversion error", http.StatusNotFound)
		log.Warn("Сonversion error")
		return
	}
	if err := ur.GetUserWithID(ctx, w, id, log); err != nil {
		log.Warn("Failed to get user tasks", "user_id", id, "error", err)
	}
}
func main()
