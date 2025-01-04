package service

import (
	"context"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
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

func (s *Service) GetOneWord(ctx context.Context, username string) (*model.Word, error) {
	words, err := s.repo.GetWordsForNotification(ctx, username, 1)
	if err != nil {
		return nil, err
	}
	if len(words) == 0 {
		return nil, nil
	}
	return words[0], nil
}

func (s *Service) AddCorrectAnswerToWord(ctx context.Context, id int) error {
	return s.repo.AddCorrectAnswerToWord(ctx, id)
}

func (s *Service) GetWordByID(ctx context.Context, id int) (*model.Word, error) {
	return s.repo.GetWordByID(ctx, id)
}
