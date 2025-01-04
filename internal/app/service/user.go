package service

import (
	"context"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) CreateUser(ctx context.Context, user *model.User) error {
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) SetWordsCount(ctx context.Context, username string, count int) error {
	return s.repo.SetWordsCount(ctx, username, count)
}
