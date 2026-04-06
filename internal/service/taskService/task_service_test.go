package task_test

import (
	"context"
	"net/url"
	"testing"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/ArsHighway/Tasks-PSQL/internal/service/mocks"
	task "github.com/ArsHighway/Tasks-PSQL/internal/service/taskService"
)

func TestTaskService_GetTaskWithID(t *testing.T) {
	t.Parallel()

	mockRepo := &mocks.MockTaskRepo{
		TaskToReturn: &models.Task{ID: 1, Title: "Task 1"},
	}

	svc := task.NewTaskService(mockRepo)
	got, err := svc.GetTaskWithID(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatalf("expected task, got nil")
	}
	if got.ID != 1 {
		t.Fatalf("expected ID=1, got %d", got.ID)
	}
	if mockRepo.ReceivedID != 1 {
		t.Fatalf("expected repo called with id=1, got %d", mockRepo.ReceivedID)
	}
}

func TestTaskService_GetTasks(t *testing.T) {
	t.Parallel()

	mockRepo := &mocks.MockTaskRepo{
		TasksToReturn: []models.Task{
			{ID: 1, Title: "Task 1", Status: "pending"},
			{ID: 2, Title: "Task 2", Status: "pending"},
		},
	}

	svc := task.NewTaskService(mockRepo)
	params := url.Values{}
	params.Set("status", "pending")
	params.Set("limit", "10")
	params.Set("page", "1")

	got, err := svc.GetTasks(context.Background(), params)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 tasks, got %d", len(got))
	}
	if len(mockRepo.ReceivedArgs) != 1 || mockRepo.ReceivedArgs[0] != "pending" {
		t.Fatalf("expected args=[pending], got %#v", mockRepo.ReceivedArgs)
	}
	if mockRepo.ReceivedBaseQuery == "" {
		t.Fatalf("expected baseQuery to be set")
	}
}

func TestTaskService_CreateTask(t *testing.T) {
	t.Parallel()
	mockRepo := &mocks.MockTaskRepo{
		TaskToReturn: &models.Task{ID: 1, Title: "Task 1"},
	}
	svc := task.NewTaskService(mockRepo)
	in := &models.Task{ID: 1, Title: "Task 1"}
	got, err := svc.CreateTask(context.Background(), in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID == 0 {
		t.Fatalf("expected ID to be set")
	}
	if got.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set")
	}
	if mockRepo.ReceivedTask == nil {
		t.Fatalf("expected repo CreateTask to be called with task")
	}
	if mockRepo.ReceivedTask.CreatedAt.IsZero() {
		t.Fatalf("expected CreatedAt to be set")
	}
}

func TestTaskService_DeleteTask(t *testing.T) {
	t.Parallel()
	mockRepo := &mocks.MockTaskRepo{
		TaskToReturn: &models.Task{ID: 1},
	}
	svc := task.NewTaskService(mockRepo)
	id := 1
	err := svc.DeleteTask(context.Background(), id)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if mockRepo.ReceivedDeletedID != id {
		t.Fatalf("expected repo called with id=%d, got %d", id, mockRepo.ReceivedDeletedID)
	}
}

func TestTaskService_GetTasksByUserID(t *testing.T) {
	t.Parallel()
	mockRepo := &mocks.MockTaskRepo{
		TasksToReturn: []models.Task{
			{ID: 1, Title: "Task 1", Status: "pending", UserID: 1},
			{ID: 2, Title: "Task 2", Status: "pending", UserID: 1},
		},
	}
	svc := task.NewTaskService(mockRepo)
	userID := 1
	tasks, err := svc.GetTasksByUserID(context.Background(), userID)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatalf("expected tasks, got empty slice")
	}
	for _, v := range tasks {
		if v.UserID != userID {
			t.Fatalf("the user IDs don't match")
		}
	}
	if mockRepo.ReceivedUserID != userID {
		t.Fatalf("repo called with wrong userID")
	}
}

func TestTaskService_UpdateTask(t *testing.T) {
	t.Parallel()
	mockRepo := &mocks.MockTaskRepo{
		TaskToReturn: &models.Task{ID: 1, Title: "Task 1"},
	}
	svc := task.NewTaskService(mockRepo)
	id := 1
	in := &models.Task{ID: 1, Title: "Task 1"}
	got, err := svc.UpdateTask(context.Background(), id, in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID == 0 {
		t.Fatalf("expected ID to be set")
	}
	if mockRepo.ReceivedID != id {
		t.Fatalf("expected repo called with id=%d, got %d", id, mockRepo.ReceivedID)
	}
	if mockRepo.ReceivedTask == nil {
		t.Fatalf("expected repo called with task")
	}
}

func TestTaskService_PatchTask(t *testing.T) {
	t.Parallel()
	mockRepo := &mocks.MockTaskRepo{
		TaskToReturn: &models.Task{ID: 1},
		ErrToReturn:  nil,
	}
	svc := task.NewTaskService(mockRepo)
	id := 1
	updates := map[string]interface{}{
		"title": "Task 2",
	}
	got, err := svc.PatchTask(context.Background(), id, updates)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID == 0 {
		t.Fatalf("expected ID to be set")
	}
	if mockRepo.ReceivedID != id {
		t.Fatalf("expected repo called with id=%d, got %d", id, mockRepo.ReceivedID)
	}
	if mockRepo.ReceivedUpdates == nil {
		t.Fatalf("expected updates to be set")
	}
	if mockRepo.ReceivedUpdates["title"] != "Task 2" {
		t.Fatalf("expected updates to be set to Task 2")
	}
}
