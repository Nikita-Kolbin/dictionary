package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
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
		logger.Error(ctx, "can't, get updates", "err", err)
	}

	for _, u := range updates {
		if u.Message != nil {
			logger.Info(ctx, "fetch message", "text", u.Message.Text, "from", u.Message.From.Username)
			s.processCommand(ctx, u.Message)
		} else if u.CallbackQuery != nil {
			s.processCallback(ctx, u.CallbackQuery)
		}
	}

	if len(updates) > 0 {
		s.tgOffset = updates[len(updates)-1].ID + 1
	}
}

func (s *Service) processCommand(ctx context.Context, msg *model.Message) {
	chatID := msg.Chat.ID

	if msg.Text == "" || msg.Text[0] != '/' {
		_, err := s.tgClient.Send(chatID, model.UnknownCommandMSG, false)
		if err != nil {
			logger.Error(ctx, "can't, send message", "err", err)
		}
		return
	}

	command, arg, _ := strings.Cut(msg.Text, " ")
	_ = arg

	var err error
	switch command {

	case model.HelpCMD:
		_, err = s.tgClient.Send(chatID, model.HelpMSG, false)
	case model.StartCMD:
		_, err = s.tgClient.Send(chatID, s.createUserTG(ctx, msg), false)

	case model.AddCMD:
		_, err = s.tgClient.Send(chatID, s.addWordTG(ctx, msg, arg), false)
	case model.GetCMD:
		text, wordID := s.getOneWordTG(ctx, msg.From.Username)
		err = s.SendWithKeyboard(text, wordID, chatID)

	case model.AddTimeCMD:
		_, err = s.tgClient.Send(chatID, s.addNotificationTimeTG(ctx, msg, arg), false)
	case model.GetTimeCMD:
		_, err = s.tgClient.Send(chatID, s.getNotificationTimeTG(ctx, msg), false)
	case model.DelTimeCMD:
		_, err = s.tgClient.Send(chatID, s.delNotificationTimeTG(ctx, msg, arg), false)

	case model.SetCountCMD:
		_, err = s.tgClient.Send(chatID, s.setWordCountTG(ctx, msg, arg), false)

	default:
		_, err = s.tgClient.Send(chatID, model.UnknownCommandMSG, false)
	}

	if err != nil {
		logger.Error(ctx, "can't, send message", "err", err)
	}
	return
}

func (s *Service) processCallback(ctx context.Context, cb *model.CallbackQuery) {
	data := &model.CallbackData{}
	if err := json.Unmarshal([]byte(cb.Data), data); err != nil {
		logger.Error(ctx, "can't, unmarshal callback data", "err", err)
		return
	}

	var postfix string
	if data.Correct {
		postfix = "\n" + model.GoodButton
		err := s.repo.AddCorrectAnswerToWord(ctx, data.WordID)
		if err != nil {
			logger.Error(ctx, "can't, add correct answer to word", "err", err, "id", data.WordID)
		}
	} else {
		postfix = "\n" + model.BadButton
	}

	word, err := s.repo.GetWordById(ctx, data.WordID)
	if err != nil {
		logger.Error(ctx, "can't get word by id", "err", err, "id", data.WordID)
	}
	var text string
	if word == nil {
		text = cb.Message.Text + postfix
	} else {
		text = buildWordMessage(word) + postfix
	}

	err = s.tgClient.Edit(text, data.ChatID, data.MessageID, true, nil)
	if err != nil {
		logger.Error(ctx, "can't, edit message", "err", err)
	}
}

func (s *Service) createUserTG(ctx context.Context, msg *model.Message) string {
	user := &model.User{
		Username: msg.From.Username,
		ChatID:   msg.Chat.ID,
	}
	if err := s.CreateUser(ctx, user); err != nil {
		logger.Error(ctx, "can't, create user", "err", err, "user", user.Username)
		return model.StartMSG
	}
	logger.Info(ctx, "user created", "user", user.Username)
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
		logger.Error(ctx, "can't, create word", "err", err, "word", w, "user", msg.From.Username)
		if errors.Is(err, model.ErrAlreadyExists) {
			return model.AddAlreadyExistsMSG
		}
		return model.AddErrorMSG
	}

	logger.Error(ctx, "word created", "word", w, "user", msg.From.Username)
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
		logger.Error(ctx, "can't create notification time", "err", err, "time", arg, "user", msg.From.Username)
		if errors.Is(err, model.ErrAlreadyExists) {
			return model.AddTimeAlreadyExistsMSG
		}
		if errors.Is(err, model.ErrNotificationLimit) {
			return model.AddTimeLimitMSG
		}
		return model.AddTimeErrorMSG
	}

	logger.Error(ctx, "notification time created", "time", arg, "user", msg.From.Username)
	return fmt.Sprintf(model.AddTimeSuccessMSG, arg)
}

func (s *Service) getNotificationTimeTG(ctx context.Context, msg *model.Message) string {
	times, err := s.repo.GetNotificationTimes(ctx, msg.From.Username)
	if err != nil {
		logger.Error(ctx, "can't get notification time", "err", err, "user", msg.From.Username)
		return model.GetTimeErrorMSG
	}
	if len(times) == 0 {
		return model.GetTimeEmptyMSG
	}

	logger.Error(ctx, "notification time given", "user", msg.From.Username)
	return buildNotificationTimeMessage(times)
}

func (s *Service) delNotificationTimeTG(ctx context.Context, msg *model.Message, arg string) string {
	arg = strings.TrimSpace(arg)

	t, err := parseTime(arg)
	if err != nil {
		return model.DelTimeEmptyMSG
	}

	err = s.repo.DelNotificationTime(ctx, msg.From.Username, t)
	if err != nil {
		logger.Error(ctx, "can't delete notification time", "err", err, "time", arg, "user", msg.From.Username)
		if errors.Is(err, model.ErrNotFound) {
			return model.DelTimeEmptyMSG
		}
		return model.DelTimeErrorMSG
	}

	logger.Error(ctx, "notification time deleted", "time", arg, "user", msg.From.Username)
	return fmt.Sprintf(model.DelTimeSuccessMSG, arg)
}

func (s *Service) getOneWordTG(ctx context.Context, username string) (string, int) {
	word, err := s.repo.GetWordsForNotification(ctx, username, 1)
	if err != nil {
		logger.Error(ctx, "can't get one word:", "user", username, "err", err)
		return model.GetErrorMSG, 0
	}
	if len(word) == 0 {
		return model.GetUserHaveNotWordsMSG, 0
	}
	return buildWordMessage(word[0]), word[0].ID
}

func (s *Service) setWordCountTG(ctx context.Context, msg *model.Message, arg string) string {
	arg = strings.TrimSpace(arg)

	cnt, err := strconv.Atoi(arg)
	if err != nil {
		return model.SetCountEmptyMSG
	}
	if cnt <= 0 || cnt > 25 {
		return model.SetCountEmptyMSG
	}

	err = s.repo.SetWordsCount(ctx, msg.From.Username, cnt)
	if err != nil {
		logger.Error(ctx, "can't set word count", "err", err, "count", cnt, "user", msg.From.Username)
		if errors.Is(err, model.ErrNotFound) {
			return model.SetCountUserNotFoundMSG
		}
		return model.SetCountErrorMSG
	}

	logger.Error(ctx, "word count set", "count", cnt, "user", msg.From.Username)
	return fmt.Sprintf(model.SetCountSuccessMSG, cnt)
}
