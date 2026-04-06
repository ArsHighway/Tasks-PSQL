package mocks

import (
	"context"
	"net/url"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	PatchUser(ctx context.Context, id int, updates map[string]interface{}) (*models.User, error)
	GetUserTasks(ctx context.Context, userID int) ([]models.Task, error)
	DeleteUser(ctx context.Context, id int) error
	GetUserWithID(ctx context.Context, id int) (*models.User, error)
}

type MockUserServ struct {
	UserToReturn *models.User
	ErrToReturn  error

	ReceivedID      int
	ReceivedParams  url.Values
	ReceivedUser    *models.User
	ReceivedUpdates map[string]interface{}
}

func (m *MockUserServ) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	m.ReceivedUser = u
	return m.UserToReturn, m.ErrToReturn
}

func (m *MockUserServ) GetUserWithID(ctx context.Context, id int) (*models.User, error) {
	m.ReceivedID = id
	return m.UserToReturn, m.ErrToReturn
}

func (m *MockUserServ) PatchUser(ctx context.Context, id int, updates map[string]interface{}) (*models.User, error) {
	m.ReceivedID = id
	m.ReceivedUpdates = updates
	return m.UserToReturn, m.ErrToReturn
}

func (m *MockUserServ) DeleteUser(ctx context.Context, id int) error {
	m.ReceivedID = id
	return nil
}

func (m *MockUserServ) GetUserTasks(ctx context.Context, userID int) ([]models.Task, error) {
	m.ReceivedID = userID
	return []models.Task{
		{ID: 1, Title: "Task 1", UserID: userID},
	}, m.ErrToReturn
}
