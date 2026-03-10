package models

import "time"

type User struct {
	ID        int       `json:"id"`         // соответствует id в таблице users
	Name      string    `json:"name"`       // name в БД
	Email     string    `json:"email"`      // email в БД
	CreatedAt time.Time `json:"created_at"` // created_at в БД
}
