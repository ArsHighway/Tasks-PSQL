package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	task "github.com/ArsHighway/Tasks-PSQL/internal/repository/taskRepository"
	"github.com/jackc/pgx/v5/pgxpool"
)

func testDBPool(t *testing.T) *pgxpool.Pool {
	t.Helper() //помечает функцию как вспомогательную, чтобы ошибки показывались в месте вызова.
	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Skip("integration: set TEST_DATABASE_URL (PostgreSQL с применёнными миграциями), например postgres://user:pass@localhost:5432/dbname")
	}
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Fatalf("pgxpool.New: %v", err)
	}
	t.Cleanup(func() { pool.Close() }) //после теста соединение автоматически закрывается
	return pool
}

func TestTaskRepository_CreateTask(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := task.NewTaskRepository(pool)

	email := fmt.Sprintf("repo_create_%d_%s@example.com", time.Now().UnixNano(), t.Name())
	var userID int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"test user", email,
	).Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	in := &models.Task{
		Title:       "title",
		Description: "desc",
		Status:      "pending",
		UserID:      userID,
		CreatedAt:   createdAt,
	}

	got, err := repo.CreateTask(ctx, in)
	if err != nil {
		t.Fatalf("CreateTask: %v", err)
	}
	if got.ID == 0 {
		t.Fatal("expected non-zero id from RETURNING")
	}

	var row models.Task
	err = pool.QueryRow(ctx,
		`SELECT id, title, description, status, user_id, created_at FROM tasks WHERE id = $1`,
		got.ID,
	).Scan(&row.ID, &row.Title, &row.Description, &row.Status, &row.UserID, &row.CreatedAt)
	if err != nil {
		t.Fatalf("verify row: %v", err)
	}
	if row.Title != in.Title || row.Description != in.Description || row.Status != in.Status || row.UserID != userID {
		t.Fatalf("stored row mismatch: %+v want title/desc/status/user_id like input", row)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, got.ID)
		if err != nil {
			t.Logf("cleanup tasks: %v", err)
		}
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
}

func TestTaskRepository_GetTaskWithID(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := task.NewTaskRepository(pool)

	var userID int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"test user", fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	createdAt := time.Now().UTC().Truncate(time.Microsecond)

	var id int
	err = pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"title", "desc", "pending", userID, createdAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed task: %v", err)
	}

	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup tasks: %v", err)
		}
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})

	got, err := repo.GetTaskWithID(ctx, id)
	if err != nil {
		t.Fatalf("GetTaskWithID: %v", err)
	}

	if got.ID != id ||
		got.Title != "title" ||
		got.Description != "desc" ||
		got.Status != "pending" ||
		got.UserID != userID {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestTaskRepository_UpdateTask(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := task.NewTaskRepository(pool)
	var userID int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"test user", fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	createdAt := time.Now().UTC().Truncate(time.Microsecond)

	var id int
	err = pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"title", "desc", "pending", userID, createdAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed task: %v", err)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup tasks: %v", err)
		}
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
	in := models.Task{
		Title:       "title",
		Description: "desc",
		Status:      "pending",
	}
	got, err := repo.UpdateTask(ctx, id, &in)
	if err != nil {
		t.Fatalf("UpdateTask: %v", err)
	}
	if got.ID != id ||
		got.Title != in.Title ||
		got.Description != in.Description ||
		got.Status != in.Status ||
		got.UserID != userID {
		t.Fatalf("unexpected task: %+v", got)
	}
}

func TestTaskRepository_PatchTask(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := task.NewTaskRepository(pool)
	var userID int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"test user", fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	createdAt := time.Now().UTC().Truncate(time.Microsecond)

	var id int
	err = pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"title", "desc", "pending", userID, createdAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed task: %v", err)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup tasks: %v", err)
		}
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
	updates := map[string]interface{}{
		"title":       "новый заголовок",
		"description": "обновлённое описание",
		"status":      "done",
	}
	parts := []string{"title", "description", "status"}
	arg := []interface{}{"новый заголовок", "обновлённое описание", "done"}
	got, err := repo.PatchTask(ctx, id, updates, parts, arg)
	if err != nil {
		t.Fatalf("PatchTask: %v", err)
	}
	if got.ID != id ||
		got.Title != "новый заголовок" ||
		got.Description != "обновлённое описание" ||
		got.Status != "done" ||
		got.UserID != userID {
		t.Fatalf("unexpected task after patch: %+v", got)
	}
}

func TestTaskRepository_DeleteTask(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := task.NewTaskRepository(pool)
	var userID int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"test user", fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	createdAt := time.Now().UTC().Truncate(time.Microsecond)

	var id int
	err = pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"title", "desc", "pending", userID, createdAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed task: %v", err)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup tasks: %v", err)
		}
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
	err = repo.DeleteTask(ctx, id)
	if err != nil {
		t.Fatalf("DeleteTask: %v", err)
	}
}

func TestTaskRepository_GetTasksByUserID(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := task.NewTaskRepository(pool)
	var userID int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"test user", fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	createdAt := time.Now().UTC().Truncate(time.Microsecond)

	var id int
	err = pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"title", "desc", "pending", userID, createdAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed task: %v", err)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup tasks: %v", err)
		}
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})

	tasks, err := repo.GetTasksByUserID(ctx, userID)
	if err != nil {
		t.Fatalf("GetTasksByUserID: %v", err)
	}
	got := tasks[0]
	if got.Title != "title" || got.Description != "desc" || got.Status != "pending" {
		t.Errorf("unexpected task: %+v", got)
	}
	if got.UserID != userID {
		t.Errorf("unexpected userID: got %d, want %d", got.UserID, userID)
	}
}

func TestTaskRepository_GetTasks(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := task.NewTaskRepository(pool)
	var userID int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"test user", fmt.Sprintf("test_%d@example.com", time.Now().UnixNano()),
	).Scan(&userID)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}
	createdAt := time.Now().UTC().Truncate(time.Microsecond)
	var id int
	err = pool.QueryRow(ctx,
		`INSERT INTO tasks (title, description, status, user_id, created_at)
		 VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		"title", "desc", "pending", userID, createdAt,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed task: %v", err)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM tasks WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup tasks: %v", err)
		}
		_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
	args := []any{"pending", userID, createdAt}
	baseQuery := "SELECT * FROM tasks WHERE 1=1 AND status = $%d AND user_id = $%d AND created_at = $%d"
	tasks, err := repo.GetTasks(ctx, args, baseQuery)
	if err != nil {
		t.Fatalf("GetTasks: %v", err)
	}
	if len(tasks) == 0 {
		t.Fatalf("expected at least 1 task, got 0")
	}
	got := tasks[0]
	if got.CreatedAt.Equal(createdAt) || got.Status != "pending" {
		t.Errorf("unexpected task: %+v", got)
	}
}
