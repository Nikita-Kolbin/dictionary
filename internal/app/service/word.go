package service

import (
	"context"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) CreateWord(ctx context.Context, word *model.Word) error {
	return s.repo.CreateWord(ctx, word)
}

func (s *Service) GetWordsForNotification(ctx context.Context, username string) ([]*model.Word, error) {
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		logger.Error(ctx, "can't get user", "err", err, "username", username)
	}

	limit := 10
	if user != nil {
		limit = user.NotificationWordCount
	}

	return s.repo.GetWordsForNotification(ctx, username, limit)
}
