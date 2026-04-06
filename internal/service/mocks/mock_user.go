package mocks

import (
	"context"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
)

type MockUserRepo struct {
	UserToReturn *models.User
	ErrToReturn  error

	ReceivedUser    *models.User
	ReceivedID      int
	ReceivedUpdates map[string]interface{}
	ReceivedParts   []string
	ReceivedArg     []interface{}
}

func (m *MockUserRepo) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	m.ReceivedUser = u
	return m.UserToReturn, m.ErrToReturn
}

func (m *MockUserRepo) GetUserWithID(ctx context.Context, id int) (*models.User, error) {
	m.ReceivedID = id
	return m.UserToReturn, m.ErrToReturn
}

func (m *MockUserRepo) PatchUser(ctx context.Context, id int, updates map[string]interface{}, parts []string, arg []interface{}) (*models.User, error) {
	m.ReceivedID = id
	m.ReceivedUpdates = updates
	m.ReceivedParts = parts
	m.ReceivedArg = arg
	if m.UserToReturn != nil && m.UserToReturn.CreatedAt.IsZero() {
		m.UserToReturn.CreatedAt = time.Now()
	}
	if m.UserToReturn == nil {
		return nil, m.ErrToReturn
	}
	return m.UserToReturn, m.ErrToReturn
}

func (m *MockUserRepo) DeleteUser(ctx context.Context, id int) error {
	m.ReceivedID = id
	return nil
}
