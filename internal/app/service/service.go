package service

import (
	"context"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

type repository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUser(ctx context.Context, username string) (*model.User, error)
	GetUsers(ctx context.Context, usernames []string) ([]*model.User, error)
	SetWordsCount(ctx context.Context, username string, count int) error

	CreateWord(ctx context.Context, word *model.Word) error
	GetWordByID(ctx context.Context, id int) (*model.Word, error)
	DeleteWord(ctx context.Context, word, username string) error
	GetWordsForNotification(ctx context.Context, username string, limit int) ([]*model.Word, error)
	AddCorrectAnswerToWord(ctx context.Context, id int) error
	GetAllUserWords(ctx context.Context, username string) ([]*model.Word, error)
	UpdateUserLastBackup(ctx context.Context, username string) error

	GetNotificationTimes(ctx context.Context, username string) ([]time.Time, error)
	AddNotificationTime(ctx context.Context, username string, t time.Time) error
	DelNotificationTime(ctx context.Context, username string, t time.Time) error
	GetUsernamesByTime(ctx context.Context, t time.Time) ([]string, error)
}

type tgClient interface {
	Updates(offset, limit int) ([]*model.Update, error)
	Send(chatID int, msg string, withFormat bool) (*model.Response, error)
	Edit(msg string, chatID, msgID int, withFormat bool, key *model.InlineKeyboardMarkup) error
	SendDocument(chatID int, filePath string) (*model.Response, error)
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
