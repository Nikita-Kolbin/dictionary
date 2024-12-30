package model

import "time"

type User struct {
	Username              string    `db:"username"`
	ChatID                int       `db:"chat_id"`
	NotificationWordCount int       `db:"notification_word_count"`
	Created               time.Time `db:"created"`
}
