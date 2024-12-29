package model

import "errors"

var ErrAlreadyExists = errors.New("already exists")

const PostgresUniqueConstraint = "23505"

const TelegramUpdateLimit = 100
