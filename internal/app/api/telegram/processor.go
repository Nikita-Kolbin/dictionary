package telegram

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) RunTelegramProcessor(ctx context.Context) {
	go func() {
		for {
			t.processUpdates(ctx)
			time.Sleep(time.Millisecond)
		}
	}()
}

func (t *Telegram) processUpdates(ctx context.Context) {
	updates, err := t.srv.Updates()
	if err != nil {
		logger.Error(ctx, "can't, get updates", "err", err)
	}

	for _, u := range updates {
		if u.Message != nil {
			logger.Info(ctx, "fetch message", "text", u.Message.Text, "from", u.Message.From.Username)
			t.processCommand(ctx, u.Message)
		} else if u.CallbackQuery != nil {
			logger.Info(ctx, "fetch callback", "from", u.CallbackQuery.From, "data", u.CallbackQuery.Data)
			t.processCallback(ctx, u.CallbackQuery)
		}
	}
}

func (t *Telegram) processCommand(ctx context.Context, msg *model.Message) { //nolint
	chatID := msg.Chat.ID

	if msg.Text == "" || msg.Text[0] != '/' {
		_, err := t.srv.Send(chatID, model.UnknownCommandMSG, false)
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
		_, err = t.srv.Send(chatID, model.HelpMSG, false)
	case model.StartCMD:
		text := t.createUserTG(ctx, msg)
		_, err = t.srv.Send(chatID, text, false)

	case model.AddCMD:
		_, err = t.srv.Send(chatID, t.addWordTG(ctx, msg, arg), false)
	case model.GetCMD:
		text, wordID := t.getOneWordTG(ctx, msg.From.Username)
		err = t.srv.SendWithKeyboard(text, wordID, chatID)
	case model.DelCMD:
		text := t.delWordTG(ctx, msg, arg)
		_, err = t.srv.Send(chatID, text, false)

	case model.AddTimeCMD:
		text := t.addNotificationTimeTG(ctx, msg, arg)
		_, err = t.srv.Send(chatID, text, false)
	case model.GetTimeCMD:
		text := t.getNotificationTimeTG(ctx, msg)
		_, err = t.srv.Send(chatID, text, false)
	case model.DelTimeCMD:
		text := t.delNotificationTimeTG(ctx, msg, arg)
		_, err = t.srv.Send(chatID, text, false)

	case model.SetCountCMD:
		_, err = t.srv.Send(chatID, t.setWordCountTG(ctx, msg, arg), false)

	case model.BackupCMD:
		var filePath, text string
		filePath, text = t.backupTG(ctx, msg)
		if len(filePath) > 0 {
			_, err = t.srv.SendWithDocument(chatID, filePath)
		} else {
			_, err = t.srv.Send(chatID, text, false)
		}

	default:
		_, err = t.srv.Send(chatID, model.UnknownCommandMSG, false)
	}

	if err != nil {
		logger.Error(ctx, "can't, send message", "err", err)
	}
}

func (t *Telegram) processCallback(ctx context.Context, cb *model.CallbackQuery) {
	data := &model.CallbackData{}
	if err := json.Unmarshal([]byte(cb.Data), data); err != nil {
		logger.Error(ctx, "can't, unmarshal callback data", "err", err)
		return
	}

	var postfix string
	if data.Correct {
		postfix = "\n" + model.GoodButton
		err := t.srv.AddCorrectAnswerToWord(ctx, data.WordID)
		if err != nil {
			logger.Error(ctx, "can't, add correct answer to word", "err", err, "id", data.WordID)
		}
	} else {
		postfix = "\n" + model.BadButton
	}

	word, err := t.srv.GetWordByID(ctx, data.WordID)
	if err != nil {
		logger.Error(ctx, "can't get word by id", "err", err, "id", data.WordID)
	}
	var text string
	if word == nil {
		text = cb.Message.Text + postfix
	} else {
		text = t.srv.BuildWordMessage(word) + postfix
	}

	err = t.srv.Edit(text, data.ChatID, data.MessageID, true, nil)
	if err != nil {
		logger.Error(ctx, "can't, edit message", "err", err)
	}
}
