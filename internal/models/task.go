package models

import "time"

type Task struct {
	ID          int       `json:"id"`          // id задачи
	Title       string    `json:"title"`       // название задачи
	Description string    `json:"description"` // описание задачи
	Status      string    `json:"status"`      // статус: pending, done и т.д.
	UserID      int       `json:"user_id"`     // id пользователя, которому принадлежит задача
	CreatedAt   time.Time `json:"created_at"`  // дата создания
}
