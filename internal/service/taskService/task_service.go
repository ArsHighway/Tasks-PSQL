package task

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/errs"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	task "github.com/ArsHighway/Tasks-PSQL/internal/repository/taskRepository"
	"github.com/jackc/pgx/v5"
)

type taskService struct {
	repo task.TaskRepository
}

type TaskService interface {
	GetTaskWithID(ctx context.Context, id int) (*models.Task, error)
	GetTasks(ctx context.Context, params url.Values) ([]models.Task, error)
	GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error)
	CreateTask(ctx context.Context, t *models.Task) (*models.Task, error)
	UpdateTask(ctx context.Context, id int, t *models.Task) (*models.Task, error)
	PatchTask(ctx context.Context, id int, updates map[string]interface{}) (*models.Task, error)
	DeleteTask(ctx context.Context, id int) error
}

func NewTaskService(repo task.TaskRepository) *taskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	if t.CreatedAt.IsZero() {
		t.CreatedAt = time.Now()
	}
	task, err := s.repo.CreateTask(ctx, t)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *taskService) GetTaskWithID(ctx context.Context, id int) (*models.Task, error) {
	task, err := s.repo.GetTaskWithID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, pgx.ErrNoRows
		}
		return nil, err
	}
	return task, nil
}

func (s *taskService) UpdateTask(ctx context.Context, id int, t *models.Task) (*models.Task, error) {
	task, err := s.repo.UpdateTask(ctx, id, t)
	if err != nil {
		if errors.Is(err, errs.ErrTaskNotFound) {
			return nil, errs.ErrTaskNotFound
		}
		if errors.Is(err, errs.ErrInvalidTask) {
			return nil, errs.ErrInvalidTask
		}
		return nil, err
	}
	return task, err
}

func (s *taskService) PatchTask(ctx context.Context, id int, updates map[string]interface{}) (*models.Task, error) {
	var arg []interface{}
	c := 1
	allowed := map[string]bool{
		"title":       true,
		"description": true,
		"status":      true,
	}
	parts := []string{}
	for k, v := range updates {
		if !allowed[k] {
			continue
		}
		arg = append(arg, v)
		parts = append(parts, fmt.Sprintf("%v =$%d", k, c))
		c++
	}
	task, err := s.repo.PatchTask(ctx, id, updates, parts, arg)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidTask) {
			return nil, errs.ErrInvalidTask
		}
		if errors.Is(err, errs.ErrTaskNotFound) {
			return nil, errs.ErrTaskNotFound
		}
		return nil, err
	}
	return task, nil
}

func (s *taskService) DeleteTask(ctx context.Context, id int) error {
	err := s.repo.DeleteTask(ctx, id)
	if errors.Is(err, errs.ErrTaskNotFound) {
		return errs.ErrTaskNotFound
	}
	return nil
}

func (s *taskService) GetTasks(ctx context.Context, params url.Values) ([]models.Task, error) {
	baseQuery := "SELECT * FROM tasks WHERE 1=1"
	args := []any{}
	if status := params.Get("status"); status != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", len(args)+1)
		args = append(args, status)
	}
	if userID := params.Get("user_id"); userID != "" {
		baseQuery += fmt.Sprintf(" AND user_id = $%d", len(args)+1)
		args = append(args, userID)
	}
	if createAt := params.Get("created_at"); createAt != "" {
		baseQuery += fmt.Sprintf(" AND created_at = $%d", len(args)+1)
		args = append(args, createAt)
	}
	allowed := map[string]bool{
		"status":     true,
		"user_id":    true,
		"created_at": true,
	}
	sortBy := params.Get("sort_by")
	orderBy := params.Get("order")
	if sortBy != "" && allowed[sortBy] {
		if orderBy != "desc" {
			orderBy = "asc"
		}
		baseQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, orderBy)
	} else {
		baseQuery += " ORDER BY created_at DESC"
	}

	limit := 10
	page := 1
	if l := params.Get("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 100 {
			limit = parsed
		} else {
			return nil, errs.ErrBadConvertation
		}
	}
	if p := params.Get("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		} else {
			return nil, errs.ErrBadConvertation
		}
	}
	baseQuery += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, (page-1)*limit)
	tasks, err := s.repo.GetTasks(ctx, args, baseQuery)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidTask) {
			return nil, errs.ErrInvalidTask
		}
		if errors.Is(err, errs.ErrTaskNotFound) {
			return nil, errs.ErrTaskNotFound
		}
		if errors.Is(err, errs.ErrNotValidFields) {
			return nil, errs.ErrNotValidFields
		}
		return nil, err
	}
	return tasks, nil
}

func (s *taskService) GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error) {
	tasks, err := s.repo.GetTasksByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, errs.ErrInvalidTask) {
			return nil, errs.ErrInvalidTask
		}
		if errors.Is(err, errs.ErrTaskNotFound) {
			return nil, errs.ErrTaskNotFound
		}
		return nil, err
	}
	return tasks, nil
}
