package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	log := slog.With("handler", "CreateTask", "method", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}

	var t repository.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Warn("JSON decoding failed", "error", err)
		return
	}

	task, err := h.repo.CreateTask(ctx, &t, *log)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		log.Error("Create task failed", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(task); err != nil {
		log.Warn("JSON encoding failed", "error", err)
		return
	}

	log.Info("Task created successfully", "taskID", task.ID)
}

func (h *TaskHandler) GetTaskWithID(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "GetTaskWithID",
		"request_method", r.Method)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Сonversion error", http.StatusNotFound)
		log.Warn("Сonversion error", "error", err)
		return
	}
	t, err := h.repo.GetTaskWithID(ctx, w, id, *log)
	if err != nil {
		http.Error(w, "Failed to get Task", http.StatusInternalServerError)
		log.Warn("Failed to get Task")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Warn("Failed to get Task", "error", err)
		return
	}
	log.Info("task received", "task", t.Title)
}

func (h *TaskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "UpdateTask",
		"request_method", r.Method)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Сonversion error", http.StatusNotFound)
		log.Warn("Сonversion error", "error", err)
		return
	}
	var task repository.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	t, err := h.repo.UpdateTask(ctx, w, id, &task, *log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Warn("Failed to update Task", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return
	}
	log.Info("task updated", "task", t.Title)
}

func (h *TaskHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "UpdateTask",
		"request_method", r.Method)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Сonversion error", http.StatusNotFound)
		log.Warn("Сonversion error", "error", err)
		return
	}
	var updates map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updates)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	t, err := h.repo.PatchTask(ctx, w, id, updates, *log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task no found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Warn("Failed to update Task", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return
	}
	log.Info("task patch", "task", t.Title)
}

func (h *TaskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "DeleteTask",
		"request_method", r.Method)
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Сonversion error", http.StatusNotFound)
		log.Warn("Сonversion error", "error", err)
		return
	}

	err = h.repo.DeleteTask(ctx, w, id, *log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task no found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Warn("Failed to update Task", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]interface{}{
		"message": "Task deleted successfully",
		"taskID":  id,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Warn("JSON encoding failed", "error", err)
	}

	log.Info("Task delete", "task", id)
}

// func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
// 	log := slog.With("handler", "TaskHandler", "method", r.Method)
// 	ctx, cancle := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancle()
// 	if r.Method != http.MethodGet {
// 		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
// 		log.Warn("Method not allowed")
// 		return
// 	}
// 	status := r.URL.Query().Get("status")
// 	t, err := h.repo.GetTasks(ctx, status, *log)
// 	if err != nil {
// 		if errors.Is(err, pgx.ErrNoRows) {
// 			http.Error(w, "Task no found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
// 		}
// 		log.Warn("Failed to update Task", "error", err)
// 		return
// 	}
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	err = json.NewEncoder(w).Encode(t)
// 	if err != nil {
// 		log.Warn("JSON encoding failed", "error", err)
// 		return
// 	}
// 	log.Info("task received", "status", status)

// }

func (h *TaskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "GetTasks", "method", r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	ctx, cancle := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancle()
	params := r.URL.Query()
	baseQuery := "SELECT * FROM tasks WHERE 1=1"
	args := []any{}
	if status := params.Get("status"); status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, status)
	}
	if userID := params.Get("user_id"); userID != "" {
		baseQuery += fmt.Sprintf(" AND user_id = $%d", len(args)+1)
		args = append(args, userID)
	}
	if createAt := params.Get("created_at"); createAt != "" {
		baseQuery += fmt.Sprintf(" AND created_at = $%d", len(args)+1)
		args = append(args, createAt)
	}
	allowed := map[string]bool{
		"status":     true,
		"user_id":    true,
		"created_at": true,
	}
	sortBy := params.Get("sort_by")
	orderBy := params.Get("order")
	if sortBy != "" && allowed[sortBy] {
		if orderBy != "desc" {
			orderBy = "asc"
		}
		baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, orderBy)
	} else {
		baseQuery += " ORDER BY created_at DESC"
		log.Warn("not allowed or parametr is null")
	}

	limit := 10
	page := 1
	if l := params.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		} else {
			http.Error(w, "Problem to convertation", http.StatusBadRequest)
			log.Warn("Problem to convertation")
		}
	}
	if p := params.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		} else {
			http.Error(w, "Problem to convertation", http.StatusBadRequest)
			log.Warn("Problem to convertation")
		}
	}
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit)

	t, err := h.repo.GetTasks(ctx, args, baseQuery, log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task no found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Warn("Failed to update Task", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(t)
	if err != nil {
		log.Warn("JSON encoding failed", "error", err)
		return
	}
	log.Info("tasks received")
}
