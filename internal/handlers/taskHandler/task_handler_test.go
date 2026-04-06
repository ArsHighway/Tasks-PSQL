package task_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ArsHighway/Tasks-PSQL/internal/handlers/mocks"
	task "github.com/ArsHighway/Tasks-PSQL/internal/handlers/taskHandler"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/go-chi/chi/v5"
)

func withChiURLParam(r *http.Request, name, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(name, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestTaskHandler_CreateTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task 1"},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	body := bytes.NewBufferString(`{"title":"Task 1","user_id":1}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.CreateTask(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_GetTaskWithID(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task1 "},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	req := withChiURLParam(httptest.NewRequest(http.MethodGet, "/tasks/1", nil), "id", "1")
	rec := httptest.NewRecorder()
	h.GetTaskWithID(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_GetTasks(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TasksToReturn: []models.Task{{ID: 1, Title: "Task1 "}},
		ErrToReturn:   nil,
	}
	h := task.NewTaskHandler(mockServ)
	req := httptest.NewRequest(http.MethodGet, "/tasks?status=open", nil)
	rec := httptest.NewRecorder()
	h.GetTasks(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got []models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !reflect.DeepEqual(got, mockServ.TasksToReturn) {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_UpdateTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task1 ", Description: "Description1", Status: "open"},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	body := bytes.NewBufferString(`{"title":"Task1 ","description":"Description1","status":"open","user_id":1}`)
	req := withChiURLParam(httptest.NewRequest(http.MethodPut, "/tasks/1", body), "id", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.UpdateTask(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
	if got.Description != mockServ.TaskToReturn.Description || got.Status != mockServ.TaskToReturn.Status {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_PatchTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		TaskToReturn: &models.Task{ID: 1, Title: "Task1 ", Description: "Description1", Status: "open"},
		ErrToReturn:  nil,
	}
	h := task.NewTaskHandler(mockServ)
	body := bytes.NewBufferString(`{"title":"Task1 ","description":"Description1","status":"open","user_id":1}`)
	req := withChiURLParam(httptest.NewRequest(http.MethodPatch, "/tasks/1", body), "id", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.PatchTask(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.Task
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.ID != mockServ.TaskToReturn.ID || got.Title != mockServ.TaskToReturn.Title {
		t.Fatalf("unexpected body: %+v", got)
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockTaskServ{
		ErrToReturn: nil,
	}
	h := task.NewTaskHandler(mockServ)
	req := withChiURLParam(httptest.NewRequest(http.MethodDelete, "/tasks/1", nil), "id", "1")
	rec := httptest.NewRecorder()
	h.DeleteTask(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got struct {
		Message string `json:"message"`
		TaskID  int    `json:"taskID"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Message != "Task deleted successfully" || got.TaskID != 1 {
		t.Fatalf("unexpected body: %+v", got)
	}
}
