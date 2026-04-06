package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/errs"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	task "github.com/ArsHighway/Tasks-PSQL/internal/repository/taskRepository"
	user "github.com/ArsHighway/Tasks-PSQL/internal/repository/userRepository"
	"github.com/jackc/pgx/v5"
)

type userService struct {
	userRepo user.UserRepository
	taskRepo task.TaskRepository
}

func NewUserService(userRepo user.UserRepository, taskRepo task.TaskRepository) *userService {
	return &userService{userRepo: userRepo, taskRepo: taskRepo}
}

type UserService interface {
	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	GetUserWithID(ctx context.Context, id int) (*models.User, error)
	PatchUser(ctx context.Context, id int, updates map[string]interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
	GetUserTasks(ctx context.Context, userID int) ([]models.Task, error)
}

func (s *userService) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	return s.userRepo.CreateUser(ctx, u)
}

func (s *userService) GetUserWithID(ctx context.Context, id int) (*models.User, error) {
	user, err := s.userRepo.GetUserWithID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return user, nil
}

func (s *userService) PatchUser(ctx context.Context, id int, updates map[string]interface{}) (*models.User, error) {
	var arg []interface{}
	c := 1
	allowed := map[string]bool{
		"name":  true,
		"email": true,
	}
	parts := []string{}
	for k, v := range updates {
		if !allowed[k] {
			continue
		}
		arg = append(arg, v)
		parts = append(parts, fmt.Sprintf("%s = $%d", k, c))
		c++
	}
	return s.userRepo.PatchUser(ctx, id, updates, parts, arg)
}

func (s *userService) DeleteUser(ctx context.Context, id int) error {
	err := s.userRepo.DeleteUser(ctx, id)
	if errors.Is(err, errs.ErrUserNotFound) {
		return errs.ErrUserNotFound
	}
	return err
}

func (s *userService) GetUserTasks(ctx context.Context, userID int) ([]models.Task, error) {
	return s.taskRepo.GetTasksByUserID(ctx, userID)
}
