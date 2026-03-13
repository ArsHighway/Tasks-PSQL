package routers

import (
	"github.com/ArsHighway/Tasks-PSQL/handlers"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(userHandler *handlers.UserHandler, taskHandler *handlers.TaskHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/tasks", taskHandler.CreateTask)
	r.Get("/tasks", taskHandler.GetTasks)
	r.Get("/tasks/{id}", taskHandler.GetTaskWithID)
	r.Put("/tasks/{id}", taskHandler.UpdateTask)
	r.Patch("/tasks/{id}", taskHandler.PatchTask)
	r.Delete("/tasks/{id}", taskHandler.DeleteTask)

	r.Post("/users", userHandler.CreateUser)
	r.Get("/users/{id}", userHandler.GetUserWithID)
	r.Delete("/users/{id}", userHandler.DeleteUser)
	r.Patch("/users/{id}", userHandler.PatchUser)
	r.Get("/users/{id}/tasks", userHandler.GetTaskWithUserID)

	return r
}
