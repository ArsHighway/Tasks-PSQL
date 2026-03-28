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
	"github.com/ArsHighway/Tasks-PSQL/internal/newerr"
	"github.com/ArsHighway/Tasks-PSQL/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
)

type userHandler struct {
	serv service.UserService
}

func NewUserHandler(serv service.UserService) *userHandler {
	return &userHandler{serv: serv}
}

type UserHandler interface {
	CreateUser(w http.ResponseWriter, r *http.Request)
	GetUserWithID(w http.ResponseWriter, r *http.Request)
	PatchUser(w http.ResponseWriter, r *http.Request)
	DeleteUser(w http.ResponseWriter, r *http.Request)
	GetTaskWithUserID(w http.ResponseWriter, r *http.Request)
}

func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	log := slog.With("handler", "CreateUser", "method", r.Method)

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}

	var u models.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Creating user", "name", u.Name, "email", u.Email)
	user, err := h.serv.CreateUser(ctx, &u)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Error("Create user failed", "error", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Warn("JSON encoding failed", "error", err)
		return
	}
	log.Info("User created successfully", "userID", user.ID)
}

func (h *userHandler) GetUserWithID(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "GetUserWithID", "request_method", r.Method)
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
	log.Info("Get user", "userID", id)
	u, err := h.serv.GetUserWithID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to get user", http.StatusInternalServerError)
		}
		log.Warn("Failed to get user", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		log.Warn("JSON encoding failed", "error", err)
		return
	}
	log.Info("user received", "user", u.Name)
}

func (h *userHandler) PatchUser(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "PatchUser", "request_method", r.Method)
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
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Warn("JSON decoding failed", "error", err)
		return
	}
	log.Info("Patch user", "userID", id)
	u, err := h.serv.PatchUser(ctx, id, updates)
	if err != nil {
		switch {
		case errors.Is(err, newerr.ErrUserNotFound):
			http.Error(w, "User not found", http.StatusNotFound)
		case errors.Is(err, newerr.ErrNotValidFieldsUser):
			http.Error(w, "No valid fields to update", http.StatusBadRequest)
		default:
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Warn("Failed to patch user", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(u); err != nil {
		http.Error(w, "Problem with encode", http.StatusInternalServerError)
		log.Warn("JSON encoding failed", "error", err)
		return
	}
	log.Info("user updated", "user", u.Name)
}

func (h *userHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "DeleteUser", "request_method", r.Method)
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
	log.Info("Delete user", "userID", id)
	err = h.serv.DeleteUser(ctx, id)
	if err != nil {
		if errors.Is(err, newerr.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
		log.Warn("Failed to delete user", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]interface{}{
		"message": "User deleted successfully",
		"userID":  id,
	}
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Warn("JSON encoding failed", "error", err)
	}

	log.Info("User deleted", "userID", id)
}

func (h *userHandler) GetTaskWithUserID(w http.ResponseWriter, r *http.Request) {
	log := slog.With("handler", "GetTaskWithUserID", "request_method", r.Method)
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		log.Warn("Method not allowed")
		return
	}
	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Сonversion error", http.StatusNotFound)
		log.Warn("Сonversion error", "error", err)
		return
	}
	log.Info("Get tasks for user", "userID", id)
	tasks, err := h.serv.GetUserTasks(ctx, id)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		log.Warn("Failed to get user tasks", "error", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasks); err != nil {
		log.Warn("JSON encoding failed", "error", err)
		return
	}
	log.Info("user tasks received", "count", len(tasks))
}
