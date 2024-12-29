package service

import (
	"context"
	"log"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

const tgLimit = 100

func (s *Service) RunTelegramProcessor(ctx context.Context) {
	go func() {
		for {
			s.processUpdates(ctx)
			time.Sleep(time.Millisecond)
		}
	}()
}

func (s *Service) processUpdates(ctx context.Context) {
	updates, err := s.tgClient.Updates(s.tgOffset, tgLimit)
	if err != nil {
		log.Println("can't, get updates:", err)
	}

	for _, u := range updates {
		// TODO: Убрать ответку
		err = s.tgClient.Send(u.Message.Chat.ID, u.Message.Text)
		if err != nil {
			logger.Error(ctx, "can't, send message:", err)
		}

		logger.Info(ctx, "fetch message", "text", u.Message.Text, "from", u.Message.From.Username)
	}

	if len(updates) > 0 {
		s.tgOffset = updates[len(updates)-1].ID + 1
	}
}
