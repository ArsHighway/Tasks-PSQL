package repository_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	user "github.com/ArsHighway/Tasks-PSQL/internal/repository/userRepository"
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

func TestUserRepository_CreateUser(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := user.NewUserRepository(pool)
	email := fmt.Sprintf("repo_create_%d_%s@example.com", time.Now().UnixNano(), t.Name())
	in := &models.User{
		Name:  "User 1",
		Email: email,
	}
	got, err := repo.CreateUser(ctx, in)
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if got.ID == 0 {
		t.Fatal("expected non-zero id from RETURNING")
	}
	var row models.User
	err = pool.QueryRow(ctx,
		`SELECT id, name, email, created_at FROM users WHERE id = $1`,
		got.ID,
	).Scan(&row.ID, &row.Name, &row.Email, &row.CreatedAt)
	if err != nil {
		t.Fatalf("verify row: %v", err)
	}
	if row.Name != in.Name || row.Email != in.Email {
		t.Fatalf("stored row mismatch: %+v want %+v", row, in)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, got.ID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
}

func TestUserRepository_GetUserWithID(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := user.NewUserRepository(pool)
	email := fmt.Sprintf("repo_create_%d_%s@example.com", time.Now().UnixNano(), t.Name())
	var id int
	in := models.User{
		Name:  "User 1",
		Email: email,
	}
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		in.Name, in.Email).Scan(&id)
	if err != nil {
		t.Fatalf("verify row: %v", err)
	}
	got, err := repo.GetUserWithID(ctx, id)
	if err != nil {
		t.Fatalf("unexpected: %v", err)
	}
	if got.Name != in.Name || got.Email != in.Email {
		t.Fatalf("stored row mismatch: %+v want %+v", got, in)
	}
	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, got.ID)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
}

func TestUserRepository_DeleteUser(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := user.NewUserRepository(pool)
	email := fmt.Sprintf("repo_delete_%d_%s@example.com", time.Now().UnixNano(), t.Name())
	var id int
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		"User To Delete", email,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	if err := repo.DeleteUser(ctx, id); err != nil {
		t.Fatalf("DeleteUser: %v", err)
	}

	var exists bool
	err = pool.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`,
		id,
	).Scan(&exists)
	if err != nil {
		t.Fatalf("verify deletion: %v", err)
	}
	if exists {
		t.Fatalf("user with id %d still exists after deletion", id)
	}
}

func TestUserRepository_PatchUser(t *testing.T) {
	pool := testDBPool(t)
	ctx := context.Background()
	repo := user.NewUserRepository(pool)

	email := fmt.Sprintf("repo_patch_%d_%s@example.com", time.Now().UnixNano(), t.Name())
	var id int
	in := &models.User{
		Name:  "User 1",
		Email: email,
	}
	err := pool.QueryRow(ctx,
		`INSERT INTO users (name, email) VALUES ($1, $2) RETURNING id`,
		in.Name, in.Email,
	).Scan(&id)
	if err != nil {
		t.Fatalf("seed user: %v", err)
	}

	parts := []string{"name=$1", "email=$2"}
	args := []interface{}{"Updated Name", "updated_" + email}

	got, err := repo.PatchUser(ctx, id, nil, parts, args)
	if err != nil {
		t.Fatalf("PatchUser: %v", err)
	}

	if got.Name != "Updated Name" || got.Email != "updated_"+email {
		t.Fatalf("stored row mismatch: got %+v, want name/email updated", got)
	}

	t.Cleanup(func() {
		_, err := pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
		if err != nil {
			t.Logf("cleanup users: %v", err)
		}
	})
}
