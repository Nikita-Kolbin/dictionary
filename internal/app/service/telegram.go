package service

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (s *Service) RunTelegramProcessor(ctx context.Context) {
	go func() {
		for {
			s.processUpdates(ctx)
			time.Sleep(time.Millisecond)
		}
	}()
}

func (s *Service) processUpdates(ctx context.Context) {
	updates, err := s.tgClient.Updates(s.tgOffset, model.TelegramUpdateLimit)
	if err != nil {
		log.Println("can't, get updates:", err)
	}

	for _, u := range updates {
		logger.Info(ctx, "fetch message", "text", u.Message.Text, "from", u.Message.From.Username)

		request := s.processCommand(ctx, u.Message)
		err = s.tgClient.Send(u.Message.Chat.ID, request, true)
		if err != nil {
			logger.Error(ctx, "can't, send message:", err)
		}
	}

	if len(updates) > 0 {
		s.tgOffset = updates[len(updates)-1].ID + 1
	}
}

func (s *Service) processCommand(ctx context.Context, msg *model.Message) string {
	if msg.Text == "" || msg.Text[0] != '/' {
		return model.UnknownCommandMSG
	}

	var err error
	command, arg, _ := strings.Cut(msg.Text, " ")
	_ = arg

	switch command {
	case model.HelpCMD:
		return model.HelpMSG

	case model.StartCMD:
		user := &model.User{
			Username: msg.From.Username,
			ChatID:   msg.Chat.ID,
		}
		err = s.CreateUser(ctx, user)
		if err != nil {
			logger.Error(ctx, "can't, create user:", "err", err, "user", user.Username)
		}
		return model.StartMSG

	default:
		return model.UnknownCommandMSG
	}
}
