package user_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ArsHighway/Tasks-PSQL/internal/handlers/mocks"
	user "github.com/ArsHighway/Tasks-PSQL/internal/handlers/userHandler"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/go-chi/chi/v5"
)

func withChiURLParam(r *http.Request, name, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(name, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

func TestUserHandler_CreateUser(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	body := bytes.NewBufferString(`{"Name":"User 1","Email":"testuser@gmail.com"}`)
	req := httptest.NewRequest(http.MethodPost, "/", body)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	h.CreateUser(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
	var got models.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := mockServ.UserToReturn
	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("unexpected body: got %+v, want %+v", got, want)
	}
}

func TestUserHandler_GetUserWithID(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	req := withChiURLParam(httptest.NewRequest(http.MethodGet, "/tasks/1", nil), "id", "1")
	rec := httptest.NewRecorder()
	h.GetUserWithID(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := mockServ.UserToReturn
	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("unexpected body: got %+v, want %+v", got, want)
	}
}

func TestUserHandler_PatchUser(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	body := bytes.NewBufferString(`{"Name":"User 1","Email":"testuser@gmail.com"}`)
	req := withChiURLParam(httptest.NewRequest(http.MethodPatch, "/tasks/1", body), "id", "1")
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	h.PatchUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	var got models.User
	if err := json.NewDecoder(rec.Body).Decode(&got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	want := mockServ.UserToReturn
	if !reflect.DeepEqual(&got, want) {
		t.Fatalf("unexpected body: got %+v, want %+v", got, want)
	}
}

func TestUserHandler_DeleteUser(t *testing.T) {
	t.Parallel()
	mockServ := &mocks.MockUserServ{
		UserToReturn: &models.User{ID: 1, Name: "User 1", Email: "testuser@gmail.com"},
		ErrToReturn:  nil,
	}
	h := user.NewUserHandler(mockServ)
	req := withChiURLParam(httptest.NewRequest(http.MethodDelete, "/tasks/1", nil), "id", "1")
	rec := httptest.NewRecorder()
	h.DeleteUser(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
}
