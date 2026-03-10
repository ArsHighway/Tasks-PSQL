package routers

import (
	"github.com/ArsHighway/Tasks-PSQL/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(userHandler *handlers.UserHandler, taskHandler *handlers.TaskHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/users", userHandler.CreateUser)
	r.Post("/tasks", taskHandler.Tasks)
	r.Get("/users/{id}/tasks", userHandler.HandlerGetUserWithID)

	return r
}
