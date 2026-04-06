package mocks

import (
	"context"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
)

type MockTaskRepo struct {
	TaskToReturn  *models.Task
	TasksToReturn []models.Task
	ErrToReturn   error

	ReceivedUserID    int
	ReceivedID        int
	ReceivedArgs      []any
	ReceivedBaseQuery string
	ReceivedTask      *models.Task
	ReceivedDeletedID int
	ReceivedUpdates   map[string]interface{}
	ReceivedParts     []string
	ReceivedArg       []interface{}
}

func (m *MockTaskRepo) GetTaskWithID(ctx context.Context, id int) (*models.Task, error) {
	m.ReceivedID = id
	return m.TaskToReturn, m.ErrToReturn
}

func (m *MockTaskRepo) GetTasks(ctx context.Context, args []any, baseQuery string) ([]models.Task, error) {
	m.ReceivedArgs = args
	m.ReceivedBaseQuery = baseQuery
	return m.TasksToReturn, m.ErrToReturn
}

func (m *MockTaskRepo) CreateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	m.ReceivedTask = t
	if m.TaskToReturn != nil && m.TaskToReturn.CreatedAt.IsZero() {
		m.TaskToReturn.CreatedAt = t.CreatedAt
	}
	if m.TaskToReturn == nil {
		return t, m.ErrToReturn
	}
	return m.TaskToReturn, m.ErrToReturn
}

func (m *MockTaskRepo) UpdateTask(ctx context.Context, id int, t *models.Task) (*models.Task, error) {
	m.ReceivedID = id
	m.ReceivedTask = t
	if m.TaskToReturn != nil && m.TaskToReturn.CreatedAt.IsZero() {
		m.TaskToReturn.CreatedAt = t.CreatedAt
	}
	if m.TaskToReturn == nil {
		return t, m.ErrToReturn
	}
	return m.TaskToReturn, m.ErrToReturn
}

func (m *MockTaskRepo) PatchTask(ctx context.Context, id int, updates map[string]interface{}, parts []string, arg []interface{}) (*models.Task, error) {
	m.ReceivedID = id
	m.ReceivedUpdates = updates
	m.ReceivedParts = parts
	m.ReceivedArg = arg
	if m.TaskToReturn != nil && m.TaskToReturn.CreatedAt.IsZero() {
		m.TaskToReturn.CreatedAt = time.Now()
	}
	if m.TaskToReturn == nil {
		return nil, m.ErrToReturn
	}
	return m.TaskToReturn, m.ErrToReturn
}

func (m *MockTaskRepo) DeleteTask(ctx context.Context, id int) error {
	m.ReceivedID = id
	m.ReceivedDeletedID = id
	return m.ErrToReturn
}

func (m *MockTaskRepo) GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error) {
	m.ReceivedUserID = userID
	return m.TasksToReturn, m.ErrToReturn
}
