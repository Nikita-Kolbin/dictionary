package model

import "errors"

var (
	ErrAlreadyExists     = errors.New("already exists")
	ErrNotFound          = errors.New("not found")
	ErrNotificationLimit = errors.New("notification limit")
)

const PostgresUniqueConstraint = "23505"

const TelegramUpdateLimit = 100
