package service

import (
	"context"
	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) CreateWord(ctx context.Context, word *model.Word) error {
	return s.repo.CreateWord(ctx, word)
}
