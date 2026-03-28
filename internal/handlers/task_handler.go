package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/ArsHighway/Tasks-PSQL/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type taskHandler struct {
	serv service.TaskService
}

func NewTaskHandler(serv service.TaskService) *taskHandler {
	return &taskHandler{serv: serv}
}

type TaskHandler interface {
	CreateTask(w http.ResponseWriter, r *http.Request)
	GetTaskWithID(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	PatchTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	GetTasks(w http.ResponseWriter, r *http.Request)
}

func (h *taskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	log := slog.With("handler", "CreateTask", "method", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}

	var t models.Task

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Creating task",
		"title", t.Title,
		"user_id", t.UserID,
	)
	task, err := h.serv.CreateTask(ctx, &t)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		log.Error("Create task failed", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err = json.NewEncoder(w).Encode(task); err != nil {
		log.Warn("JSON encoding failed", "error", err)
		return
	}

	log.Info("Task created successfully", "taskID", task.ID)
}

func (h *taskHandler) GetTaskWithID(w http.ResponseWriter, r *http.Request) {
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
	log.Info("Get task", "taskID", id)
	t, err := h.serv.GetTaskWithID(ctx, id)
	if err != nil {
		http.Error(w, "Failed to get Task", http.StatusInternalServerError)
		log.Warn("Failed to get Task")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(t); err != nil {
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

func (h *taskHandler) UpdateTask(w http.ResponseWriter, r *http.Request) {
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
	var task models.Task
	err = json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Update task", "taskID", id)
	t, err := h.serv.UpdateTask(ctx, id, &task)
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
	if err = json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return
	}
	log.Info("task updated", "task", t.Title)
}

func (h *taskHandler) PatchTask(w http.ResponseWriter, r *http.Request) {
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
	log.Info("Patch task", "taskID", id)
	t, err := h.serv.PatchTask(ctx, id, updates)
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
	if err = json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return
	}
	log.Info("task patch", "task", t.Title)
}

func (h *taskHandler) DeleteTask(w http.ResponseWriter, r *http.Request) {
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
	log.Info("Update task", "taskID", id)
	err = h.serv.DeleteTask(ctx, id)
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
	if err = json.NewEncoder(w).Encode(resp); err != nil {
		log.Warn("JSON encoding failed", "error", err)
	}

	log.Info("Task delete", "task", id)
}

func (h *taskHandler) GetTasks(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "GetTasks", "method", r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	ctx, cancle := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancle()
	params := r.URL.Query()
	log.Info("Get tasks")
	t, err := h.serv.GetTasks(ctx, params)
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
