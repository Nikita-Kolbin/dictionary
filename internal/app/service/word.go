package service

import (
	"context"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) CreateWord(ctx context.Context, word *model.Word) error {
	return s.repo.CreateWord(ctx, word)
}

func (s *Service) GetWordsForNotification(ctx context.Context, username string) ([]*model.Word, error) {
	// TODO: достать из юзера
	limit := 10

	return s.repo.GetWordsForNotification(ctx, username, limit)
}
