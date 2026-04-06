package main

import (
	"log"
	"net/http"

	"github.com/ArsHighway/Tasks-PSQL/internal/config"
	taskHandler "github.com/ArsHighway/Tasks-PSQL/internal/handlers/taskHandler"
	userHandler "github.com/ArsHighway/Tasks-PSQL/internal/handlers/userHandler"
	taskRepository "github.com/ArsHighway/Tasks-PSQL/internal/repository/taskRepository"
	userRepository "github.com/ArsHighway/Tasks-PSQL/internal/repository/userRepository"
	"github.com/ArsHighway/Tasks-PSQL/internal/routers"
	taskService "github.com/ArsHighway/Tasks-PSQL/internal/service/taskService"
	userService "github.com/ArsHighway/Tasks-PSQL/internal/service/userService"
)

func main() {
	pool := config.NewPostgressPool("postgres://arsver@localhost:5432/arsver")
	defer pool.Close()

	userRepo := userRepository.NewUserRepository(pool)
	taskRepo := taskRepository.NewTaskRepository(pool)
	taskServ := taskService.NewTaskService(taskRepo)
	userServ := userService.NewUserService(userRepo, taskRepo)

	userHandler := userHandler.NewUserHandler(userServ)
	taskHandler := taskHandler.NewTaskHandler(taskServ)

	r := routers.RegisterRoutes(userHandler, taskHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
