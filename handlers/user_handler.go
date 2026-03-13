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

type UserHandler struct {
	repo *repository.UserRepository
}

func NewUserHandler(repo *repository.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	log := slog.With("handler", "CreateUsers",
		"request_method", r.Method)
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	var u repository.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Problem with decoding", http.StatusNotFound)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	if err := h.repo.CreateUser(ctx, w, &u, *log); err != nil {
		log.Warn("failed to create user", "user", u, "error", err)
		return
	}
}

func (h *UserHandler) GetTaskWithUserID(w http.ResponseWriter, r *http.Request) {
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
		log.Warn("Сonversion error", "error", err)
		return
	}
	if err := h.repo.GetTaskWithUserID(ctx, w, id, log); err != nil {
		log.Warn("Failed to get user tasks", "user_id", id, "error", err)
	}
}

func (h *UserHandler) GetUserWithID(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	log := slog.With("handler", "GetUserWithID",
		"request_method", r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Conversion error", http.StatusNotFound)
		log.Warn("Conversion error", "error", err)
		return
	}
	u, err := h.repo.GetUserWithID(ctx, id, *log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task no founded", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Warn("Failed to get user", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return
	}
	log.Info("user received", "user", u.Name)
}

func (h *UserHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	var updates map[string]interface{}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	log := slog.With("handler", "PatchUser",
		"request_method", r.Method)
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Conversion error", http.StatusNotFound)
		log.Warn("Conversion error")
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Problem with decoding", http.StatusNotFound)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	u, err := h.repo.PatchUser(ctx, id, updates, *log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task no founded", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Warn("Failed to get user", "error", err)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return
	}
	log.Info("user updated", "user", u.Name)
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	log := slog.With("handler", "DeleteUser",
		"request_method", r.Method)
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Conversion error", http.StatusNotFound)
		log.Warn("Conversion error", "error", err)
		return
	}
	err = h.repo.DeleteUser(ctx, id, *log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Task no founded", http.StatusNotFound)
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		log.Warn("Failed to get user", "error", err)
	}
	ans := map[string]interface{}{
		"message": "user deleted successfully",
		"UserID":  id,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(&ans); err != nil {
		http.Error(w, "Problem with encode", http.StatusNotFound)
		log.Warn("JSON encoding faling", "error", err)
		return
	}
	log.Info("user deleted", "userID", id)
	return
}
