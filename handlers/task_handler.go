package handlers

import (
	"context"
	"encoding/json"
	"errors"
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
	ctx, cancle := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancle()
	log := slog.With("handler", "Tasks",
		"request_method", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	var t repository.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Problem with decoding", http.StatusNotFound)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	if err := h.repo.CreateTask(ctx, w, &t, *log); err != nil {
		http.Error(w, "Failed to create Task", http.StatusInternalServerError)
		log.Warn("Failed to create Task")
		return
	}
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
