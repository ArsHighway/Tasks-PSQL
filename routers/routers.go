package routers

import (
	"github.com/ArsHighway/Tasks-PSQL/internal/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(userHandler handlers.UserHandler, taskHandler handlers.TaskHandler) chi.Router {
	r := chi.NewRouter()
	r.Route("/tasks", func(r chi.Router) {
		r.Post("/", taskHandler.CreateTask)
		r.Get("/", taskHandler.GetTasks)
		r.Get("/{id}", taskHandler.GetTaskWithID)
		r.Put("/{id}", taskHandler.UpdateTask)
		r.Patch("/{id}", taskHandler.PatchTask)
		r.Delete("/{id}", taskHandler.DeleteTask)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.CreateUser)
		r.Get("/{id}/tasks", userHandler.GetTaskWithUserID)
		r.Get("/{id}", userHandler.GetUserWithID)
		r.Delete("/{id}", userHandler.DeleteUser)
		r.Patch("/{id}", userHandler.PatchUser)
	})
	return r
}
