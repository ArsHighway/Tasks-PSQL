package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/repository"
)

type TaskHandler struct {
	repo *repository.TaskRepository
}

func NewTaskHandler(repo *repository.TaskRepository) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) Tasks(w http.ResponseWriter, r *http.Request) {
	ctx, cancle := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancle()
	log := slog.With("handler", "Tasks",
		"request_method", r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
	}
	var t repository.Task
	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, "Problem with decoding", http.StatusNotFound)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	if err := h.repo.CreateTask(ctx, w, &t, *log); err != nil {
		return
	}
}
