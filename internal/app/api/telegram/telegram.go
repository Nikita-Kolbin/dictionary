package telegram

import (
	"context"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

type service interface {
	Updates() ([]*model.Update, error)
	Send(chatID int, message string, withFormat bool) (*model.Response, error)
	SendWithKeyboard(text string, wordID, chatID int) error
	Edit(msg string, chatID, msgID int, withFormat bool, key *model.InlineKeyboardMarkup) error

	CreateUser(ctx context.Context, user *model.User) error
	SetWordsCount(ctx context.Context, username string, count int) error

	CreateWord(ctx context.Context, word *model.Word) error
	GetOneWord(ctx context.Context, username string) (*model.Word, error)
	AddCorrectAnswerToWord(ctx context.Context, id int) error
	GetWordByID(ctx context.Context, id int) (*model.Word, error)

	AddNotificationTime(ctx context.Context, username string, t time.Time) error
	GetNotificationTimes(ctx context.Context, username string) ([]time.Time, error)
	DelNotificationTime(ctx context.Context, username string, t time.Time) error
}

type Telegram struct {
	srv service
}

func New(srv service) *Telegram {
	return &Telegram{srv: srv}
}
