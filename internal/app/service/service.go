package service

import (
	"context"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

type repository interface {
	CreateUser(ctx context.Context, user *model.User) error

	CreateWord(ctx context.Context, word *model.Word) error

	GetNotificationTimes(ctx context.Context, username string) ([]time.Time, error)
	AddNotificationTime(ctx context.Context, username string, t time.Time) error
}

type tgClient interface {
	Updates(offset, limit int) ([]*model.Update, error)
	Send(chatID int, msg string, withFormat bool) error
}

type Service struct {
	tgOffset int

	repo     repository
	tgClient tgClient
}

func New(repo repository, tgCli tgClient) *Service {
	return &Service{
		repo:     repo,
		tgClient: tgCli,
	}
}
