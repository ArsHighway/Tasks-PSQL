package mocks

import (
	"context"
	"net/url"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
)

type MockTaskServ struct {
	TaskToReturn   *models.Task
	TasksToReturn  []models.Task
	TasksByUserRet []models.Task
	ErrToReturn    error

	ReceivedID      int
	ReceivedParams  url.Values
	ReceivedUserID  int
	ReceivedTask    *models.Task
	ReceivedUpdates map[string]interface{}
}

func (m *MockTaskServ) GetTaskWithID(ctx context.Context, id int) (*models.Task, error) {
	m.ReceivedID = id
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	return m.TaskToReturn, nil
}

func (m *MockTaskServ) GetTasks(ctx context.Context, params url.Values) ([]models.Task, error) {
	m.ReceivedParams = params
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	return m.TasksToReturn, nil
}

func (m *MockTaskServ) GetTasksByUserID(ctx context.Context, userID int) ([]models.Task, error) {
	m.ReceivedUserID = userID
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	return m.TasksByUserRet, nil
}

func (m *MockTaskServ) CreateTask(ctx context.Context, t *models.Task) (*models.Task, error) {
	m.ReceivedTask = t
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	return m.TaskToReturn, nil
}

func (m *MockTaskServ) UpdateTask(ctx context.Context, id int, t *models.Task) (*models.Task, error) {
	m.ReceivedID = id
	m.ReceivedTask = t
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	return m.TaskToReturn, nil
}

func (m *MockTaskServ) PatchTask(ctx context.Context, id int, updates map[string]interface{}) (*models.Task, error) {
	m.ReceivedID = id
	m.ReceivedUpdates = updates
	if m.ErrToReturn != nil {
		return nil, m.ErrToReturn
	}
	return m.TaskToReturn, nil
}

func (m *MockTaskServ) DeleteTask(ctx context.Context, id int) error {
	m.ReceivedID = id
	return m.ErrToReturn
}
