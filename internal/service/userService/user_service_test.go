package user_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/ArsHighway/Tasks-PSQL/internal/service/mocks"
	user "github.com/ArsHighway/Tasks-PSQL/internal/service/userService"
)

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()
	mockUserRepo := &mocks.MockUserRepo{
		UserToReturn: &models.User{
			ID: 1, Name: "User 1"},
	}
	mockTaskRepo := &mocks.MockTaskRepo{
		// TaskToReturn: &models.Task{ID: 1, Title: "Task 1", UserID: 1},
	}

	svc := user.NewUserService(mockUserRepo, mockTaskRepo)
	ctx := context.Background()
	in := models.User{
		ID: 1, Name: "User 1",
	}
	got, err := svc.CreateUser(ctx, &in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := mockUserRepo.UserToReturn

	if got.ID != mockUserRepo.UserToReturn.ID && got.Name != mockUserRepo.UserToReturn.Name {
		t.Fatalf("unexpected user: got %+v, want %+v", got, want)
	}
}

func TestUserService_GetUserWithID(t *testing.T) {
	t.Parallel()
	mockUserRepo := &mocks.MockUserRepo{
		UserToReturn: &models.User{
			ID: 1, Name: "User 1"},
	}
	mockTaskRepo := &mocks.MockTaskRepo{}
	svc := user.NewUserService(mockUserRepo, mockTaskRepo)
	ctx := context.Background()
	id := 1
	got, err := svc.GetUserWithID(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected task, got nil")
	}
	if got.ID != 1 {
		t.Fatalf("expected ID=1, got %d", got.ID)
	}
	if mockUserRepo.ReceivedID != 1 {
		t.Fatalf("expected repo called with id=1, got %d", mockUserRepo.ReceivedID)
	}

}

func TestUserService_PatchUser(t *testing.T) {
	t.Parallel()

	mockUserRepo := &mocks.MockUserRepo{
		UserToReturn: &models.User{
			ID: 1, Name: "User 1", Email: "test_email@gmail.com",
		},
	}
	mockTaskRepo := &mocks.MockTaskRepo{}

	svc := user.NewUserService(mockUserRepo, mockTaskRepo)
	ctx := context.Background()

	id := 1
	updates := map[string]interface{}{
		"Name":  "Updated User",
		"Email": "updated_email@gmail.com",
	}

	got, err := svc.PatchUser(ctx, id, updates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if mockUserRepo.ReceivedID != id {
		t.Errorf("expected id %d, got %d", id, mockUserRepo.ReceivedID)
	}
	if !reflect.DeepEqual(mockUserRepo.ReceivedUpdates, updates) {
		t.Errorf("expected updates %+v, got %+v", updates, mockUserRepo.ReceivedUpdates)
	}

	want := mockUserRepo.UserToReturn
	if !reflect.DeepEqual(got, want) {
		t.Errorf("unexpected user: got %+v, want %+v", got, want)
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Parallel()

	mockUserRepo := &mocks.MockUserRepo{
		UserToReturn: &models.User{
			ID: 1, Name: "User 1", Email: "test_email@gmail.com",
		},
	}
	mockTaskRepo := &mocks.MockTaskRepo{}

	svc := user.NewUserService(mockUserRepo, mockTaskRepo)
	ctx := context.Background()

	id := 1
	err := svc.DeleteUser(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
