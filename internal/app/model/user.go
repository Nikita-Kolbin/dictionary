package model

import "time"

type User struct {
	Username string    `db:"username"`
	ChatID   int       `db:"chat_id"`
	Created  time.Time `db:"created"`
}
