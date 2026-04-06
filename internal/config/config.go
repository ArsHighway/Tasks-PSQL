package config

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewPostgressPool(connect string) *pgxpool.Pool {
	pool, err := pgxpool.New(context.Background(), connect)
	if err != nil {
		log.Fatalf("Unable to connect to DB: %v", err)
	}
	return pool
}
