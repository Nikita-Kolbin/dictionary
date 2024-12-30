package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) AddNotificationTime(ctx context.Context, username string, t time.Time) error {
	times, err := s.repo.GetNotificationTimes(ctx, username)
	if err != nil {
		return fmt.Errorf("AddNotificationTime: %w", err)
	}

	if len(times) >= 3 {
		return fmt.Errorf("AddNotificationTime: %w", model.ErrNotificationLimit)
	}

	return s.repo.AddNotificationTime(ctx, username, t)
}
