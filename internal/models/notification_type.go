package models

import "time"

type NotificationType struct {
	ID         int       `db:"id"`
	Name       string    `db:"name"`        // текст уведоления
	TimeBefore int64     `db:"time_before"` // секунд до дедлайна
	Enabled    bool      `db:"enabled"`     // активность типа уведомления
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`
}
