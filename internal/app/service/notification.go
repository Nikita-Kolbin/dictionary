package service

import (
	"context"
	"fmt"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) RunNotification(ctx context.Context) {
	go func() {
		currMinute := time.Now().Minute()
		for {
			time.Sleep(time.Second)
			now := time.Now()
			nowMinute := now.Minute()
			if nowMinute == currMinute {
				continue
			}

			usernames, err := s.repo.GetUsernamesByTime(ctx, now)
			if err != nil {
				logger.Error(ctx, "can't, get usernames for notification", "err", err)
				continue
			}

			users, err := s.repo.GetUsers(ctx, usernames)
			if err != nil {
				logger.Error(ctx, "can't, get users for notification", "err", err)
				continue
			}

			for _, user := range users {
				words, err := s.repo.GetWordsForNotification(ctx, user.Username, user.NotificationWordCount)
				if err != nil {
					logger.Error(ctx, "can't, get words for notification", "err", err, "user", user)
					continue
				}
				for _, word := range words {
					text := buildWordMessage(word)
					err = s.SendWithKeyboard(text, word.ID, user.ChatID)
					if err != nil {
						logger.Error(ctx, "can't, send words for notification", "err", err, "user", user)
						continue
					}
				}
			}

			currMinute = now.Minute()
		}
	}()
}

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
