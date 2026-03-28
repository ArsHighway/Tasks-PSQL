package main

import (
	"log"
	"net/http"

	"github.com/ArsHighway/Tasks-PSQL/config"
	"github.com/ArsHighway/Tasks-PSQL/internal/handlers"
	"github.com/ArsHighway/Tasks-PSQL/internal/repository"
	"github.com/ArsHighway/Tasks-PSQL/internal/service"
	"github.com/ArsHighway/Tasks-PSQL/routers"
)

func main() {
	pool := config.NewPostgressPool("postgres://arsver@localhost:5432/arsver")
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	taskRepo := repository.NewTaskRepository(pool)
	taskServ := service.NewTaskService(taskRepo)
	userServ := service.NewUserService(userRepo, taskRepo)

	userHandler := handlers.NewUserHandler(userServ)
	taskHandler := handlers.NewTaskHandler(taskServ)

	r := routers.RegisterRoutes(userHandler, taskHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
