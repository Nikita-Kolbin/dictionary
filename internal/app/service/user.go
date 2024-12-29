package service

import (
	"context"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) CreateUser(ctx context.Context, user *model.User) error {
	return s.repo.CreateUser(ctx, user)
}
