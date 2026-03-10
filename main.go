package main

import (
	"log"
	"net/http"

	"github.com/ArsHighway/Tasks-PSQL/config"
	"github.com/ArsHighway/Tasks-PSQL/handlers"
	"github.com/ArsHighway/Tasks-PSQL/repository"
	"github.com/ArsHighway/Tasks-PSQL/routers"
)

func main() {
	pool := config.NewPostgressPool("postgres://arsver@localhost:5432/arsver")
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	taskRepo := repository.NewTaskRepository(pool)

	userHandler := handlers.NewUserHandler(userRepo)
	taskHandler := handlers.NewTaskHandler(taskRepo)

	r := routers.RegisterRoutes(userHandler, taskHandler)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
