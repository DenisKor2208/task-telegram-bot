package models

import "time"

type User struct {
	ID        int       `db:"id"`
	TgID      int       `db:"tg_id"`
	Name      string    `db:"name"`
	Timezone  string    `db:"timezone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
