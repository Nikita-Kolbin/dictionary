package service

import "github.com/Nikita-Kolbin/dictionary/internal/app/model"

type repository interface {
	// repo
}

type tgClient interface {
	Updates(offset, limit int) ([]*model.Update, error)
	Send(chatID int, msg string) error
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
