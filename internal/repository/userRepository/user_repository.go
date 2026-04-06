package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/ArsHighway/Tasks-PSQL/internal/errs"
	"github.com/ArsHighway/Tasks-PSQL/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *userRepository {
	return &userRepository{pool: pool}
}

type UserRepository interface {
	CreateUser(ctx context.Context, u *models.User) (*models.User, error)
	GetUserWithID(ctx context.Context, id int) (*models.User, error)
	PatchUser(ctx context.Context, id int, updates map[string]interface{}, parts []string, arg []interface{}) (*models.User, error)
	DeleteUser(ctx context.Context, id int) error
}

func (r *userRepository) CreateUser(ctx context.Context, u *models.User) (*models.User, error) {
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3) RETURNING id`,
		u.Name, u.Email, u.CreatedAt,
	).Scan(&u.ID)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (r *userRepository) GetUserWithID(ctx context.Context, id int) (*models.User, error) {
	var u models.User
	err := r.pool.QueryRow(ctx, `SELECT id, name, email, created_at FROM users WHERE id = $1`, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) PatchUser(ctx context.Context, id int, _ map[string]interface{}, parts []string, arg []interface{}) (*models.User, error) {
	var u models.User
	if len(parts) == 0 {
		return nil, errs.ErrNotValidFieldsUser
	}
	idx := len(arg) + 1
	sql := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(parts, ", "), idx)
	args := append(append([]interface{}{}, arg...), id)
	cmdTag, err := r.pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, errs.ErrInvalidUser
	}
	if cmdTag.RowsAffected() == 0 {
		return nil, errs.ErrUserNotFound
	}
	if err := r.pool.QueryRow(ctx, `SELECT id, name, email, created_at FROM users WHERE id = $1`, id).Scan(
		&u.ID, &u.Name, &u.Email, &u.CreatedAt,
	); err != nil {
		return nil, errs.ErrInvalidUser
	}
	return &u, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, id int) error {
	cmdTag, err := r.pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {
		return errs.ErrUserNotFound
	}
	return nil
}
