package service

import "github.com/Nikita-Kolbin/dictionary/internal/app/model"

type tgClient interface {
	Updates(offset, limit int) ([]*model.Update, error)
	Send(chatId int, msg string) error
}

type Service struct {
	tgOffset int
	tgClient tgClient
}

func New(tgCli tgClient) *Service {
	return &Service{
		tgClient: tgCli,
	}
}
