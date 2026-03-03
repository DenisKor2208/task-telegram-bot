package models

import (
	"time"
)

type Task struct {
	ID          int       `db:"id"`
	Description string    `db:"description"`
	Deadline    time.Time `db:"deadline,omitempty"`  // omitempty, если deadline может быть null
	Priority    int       `db:"priority,omitempty"`  // omitempty, если priority может быть null
	StatusID    int       `db:"status_id,omitempty"` // Ссылка на statuses, может быть null
	UserID      int       `db:"user_id"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
