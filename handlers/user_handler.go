package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/repository"
	"github.com/go-chi/chi/v5"
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
	}
	var u repository.User
	err := json.NewDecoder(r.Body).Decode(&u)
	if err != nil {
		http.Error(w, "Problem with decoding", http.StatusNotFound)
		log.Warn("JSON decoding failed", "error", err)
	}
	if err := h.repo.CreateUsers(ctx, w, &u, *log); err != nil {
		log.Warn("failed to create user", &u, "error", err)
	}
}

func (h *UserHandler) HandlerGetUserWithID(w http.ResponseWriter, r *http.Request) {
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
	if err := h.repo.GetUserWithID(ctx, w, id, log); err != nil {
		log.Warn("Failed to get user tasks", "user_id", id, "error", err)
	}
}
