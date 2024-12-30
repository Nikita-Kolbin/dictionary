package service

import (
	"context"
	"errors"
	"fmt"
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

		request, format := s.processCommand(ctx, u.Message)
		err = s.tgClient.Send(u.Message.Chat.ID, request, format)
		if err != nil {
			logger.Error(ctx, "can't, send message:", err)
		}
	}

	if len(updates) > 0 {
		s.tgOffset = updates[len(updates)-1].ID + 1
	}
}

func (s *Service) processCommand(ctx context.Context, msg *model.Message) (text string, format bool) {
	if msg.Text == "" || msg.Text[0] != '/' {
		return model.UnknownCommandMSG, false
	}

	command, arg, _ := strings.Cut(msg.Text, " ")
	_ = arg

	switch command {
	case model.HelpCMD:
		return model.HelpMSG, false
	case model.StartCMD:
		return s.createUserTG(ctx, msg), false
	case model.AddCMD:
		return s.addWordTG(ctx, msg, arg), false
	case model.AddTimeCMD:
		return s.addNotificationTimeTG(ctx, msg, arg), false
	default:
		return model.UnknownCommandMSG, false
	}
}

func (s *Service) createUserTG(ctx context.Context, msg *model.Message) string {
	user := &model.User{
		Username: msg.From.Username,
		ChatID:   msg.Chat.ID,
	}
	if err := s.CreateUser(ctx, user); err != nil {
		logger.Error(ctx, "can't, create user:", "err", err, "user", user.Username)
	}
	logger.Info(ctx, "user created:", "user", user.Username)
	return model.StartMSG
}

func (s *Service) addWordTG(ctx context.Context, msg *model.Message, arg string) string {
	w, tw, e, te := parseWord(arg)
	if len(arg) == 0 || len(w) == 0 || len(tw) == 0 {
		return model.AddEmptyMSG
	}

	word := &model.Word{
		Word:              w,
		TranslatedWord:    tw,
		Example:           e,
		TranslatedExample: te,
		Username:          msg.From.Username,
	}

	err := s.CreateWord(ctx, word)
	if err != nil {
		logger.Error(ctx, "can't, create word:", "err", err, "word", w, "user", msg.From.Username)
		if errors.Is(err, model.ErrAlreadyExists) {
			return model.AddAlreadyExistsMSG
		}
		return model.AddErrorMSG
	}

	logger.Error(ctx, "word created:", "word", w, "user", msg.From.Username)
	return fmt.Sprintf(model.AddSuccessMSG, w)
}

func (s *Service) addNotificationTimeTG(ctx context.Context, msg *model.Message, arg string) string {
	arg = strings.TrimSpace(arg)

	t, err := parseTime(arg)
	if err != nil {
		return model.AddTimeEmptyMSG
	}

	err = s.AddNotificationTime(ctx, msg.From.Username, t)
	if err != nil {
		logger.Error(ctx, "can't create notification time:", "err", err, "time", arg, "user", msg.From.Username)
		if errors.Is(err, model.ErrAlreadyExists) {
			return model.AddTimeAlreadyExistsMSG
		}
		if errors.Is(err, model.ErrNotificationLimit) {
			return model.AddTimeLimitMSG
		}
		return model.AddTimeErrorMSG
	}

	logger.Error(ctx, "notification time created:", "time", arg, "user", msg.From.Username)
	return fmt.Sprintf(model.AddTimeSuccessMSG, arg)
}
