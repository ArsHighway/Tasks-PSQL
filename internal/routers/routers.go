package routers

import (
	task "github.com/ArsHighway/Tasks-PSQL/internal/handlers/taskHandler"
	user "github.com/ArsHighway/Tasks-PSQL/internal/handlers/userHandler"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(userHandler user.UserHandler, taskHandler task.TaskHandler) chi.Router {
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
