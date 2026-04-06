package task

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/errs"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	task "github.com/ArsHighway/Tasks-PSQL/internal/repository/taskRepository"
	"github.com/Masterminds/squirrel"
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
	allowed := map[string]bool{
		"title":       true,
		"description": true,
		"status":      true,
	}
	values := make(map[string]interface{})
	for k, v := range updates {
		if allowed[k] {
			values[k] = v
		}
	}
	task, err := s.repo.PatchTask(ctx, id, values)
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

func (s *taskService) GetTasks(ctx context.Context, params url.Values) ([]models.Task, error) {
	qb := squirrel.Select("*").
		From("tasks").
		PlaceholderFormat(squirrel.Dollar)

	if status := params.Get("status"); status != "" {
		qb = qb.Where(squirrel.Eq{"status": status})
	}
	if userID := params.Get("user_id"); userID != "" {
		qb = qb.Where(squirrel.Eq{"user_id": userID})
	}
	if createdAt := params.Get("created_at"); createdAt != "" {
		qb = qb.Where(squirrel.Eq{"created_at": createdAt})
	}

	allowedSorts := map[string]string{
		"created_at": "created_at",
		"title":      "title",
		"status":     "status",
	}
	sortBy := params.Get("sort_by")
	order := strings.ToUpper(params.Get("order_by"))
	if col, ok := allowedSorts[sortBy]; ok {
		if order != "ASC" && order != "DESC" {
			order = "ASC"
		}
		qb = qb.OrderBy(fmt.Sprintf("%s %s", col, order))
	} else {
		qb = qb.OrderBy("created_at DESC")
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
	qb = qb.Limit(uint64(limit)).Offset(uint64((page - 1) * limit))

	sql, args, err := qb.ToSql()
	tasks, err := s.repo.GetTasks(ctx, args, sql)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
