package service

import (
	"context"
	"fmt"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (s *Service) RunNotification(ctx context.Context) {
	// Джоба для отправки слов по расписанию
	go s.sendMessagesJob(ctx)

	// Джоба для закрытия старых слов
	go s.autoCloseOldWordsJob(ctx)
}

func (s *Service) sendMessagesJob(ctx context.Context) {
	currMinute := time.Now().Minute()
	for {
		time.Sleep(time.Second)
		loc, _ := time.LoadLocation("Europe/Moscow")
		now := time.Now().In(loc)
		nowMinute := now.Minute()
		if nowMinute == currMinute {
			continue
		}

		// Получение юзеров с уведами на это время
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

		// Получение и рассылка слов
		for _, user := range users {
			words, err := s.repo.GetWordsForNotification(ctx, user.Username, user.NotificationWordCount)
			if err != nil {
				logger.Error(ctx, "can't, get words for notification", "err", err, "user", user)
				continue
			}
			go func(wordsCopy []*model.Word) {
				for _, word := range wordsCopy {
					time.Sleep(100 * time.Millisecond)
					text := s.BuildWordMessage(word)
					err = s.SendWithKeyboard(ctx, text, word.ID, user.ChatID, word.NeedReverseLang)
					if err != nil {
						logger.Error(ctx, "can't, send words for notification", "err", err, "user", user)
						continue
					}
				}
			}(words)
		}

		currMinute = now.Minute()
	}
}

func (s *Service) autoCloseOldWordsJob(ctx context.Context) {
	for {
		words, err := s.repo.GetOldSendWords(ctx)
		if err != nil {
			logger.Error(ctx, "can't, get old send words", "err", err)
		}

		// TODO: мб распараллелить
		for _, word := range words {
			text := s.BuildWordMessage(word)
			text += "\n" + model.BadButton

			if word.CurrentMsgID == nil {
				continue
			}

			err = s.Edit(text, word.ChatID, *word.CurrentMsgID, true, nil)
			if err != nil {
				logger.Error(ctx, "can't, edit message", "err", err)
				continue
			}

			err = s.UpdateWordCurrentMessageID(ctx, word.ID, nil)
			if err != nil {
				logger.Error(ctx, "can't, update word current message", "err", err)
				continue
			}

			logger.Info(ctx, "word auto close successful", "word_id", word.ID)
		}

		time.Sleep(time.Hour)
	}
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

func (s *Service) GetNotificationTimes(ctx context.Context, username string) ([]time.Time, error) {
	return s.repo.GetNotificationTimes(ctx, username)
}

func (s *Service) DelNotificationTime(ctx context.Context, username string, t time.Time) error {
	return s.repo.DelNotificationTime(ctx, username, t)
}
