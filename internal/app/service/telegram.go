package service

import (
	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (s *Service) Updates() ([]*model.Update, error) {
	updates, err := s.tgClient.Updates(s.tgOffset, model.TelegramUpdateLimit)

	if len(updates) > 0 {
		s.tgOffset = updates[len(updates)-1].ID + 1
	}

	return updates, err
}

func (s *Service) Send(chatID int, message string, withFormat bool) (*model.Response, error) {
	return s.tgClient.Send(chatID, message, withFormat)
}

func (s *Service) SendWithKeyboard(text string, wordID, chatID int) error {
	resp, err := s.tgClient.Send(chatID, text, true)
	if err != nil {
		return err
	}

	if wordID > 0 && chatID > 0 && resp.Result.MessageID > 0 {
		key := getKeyTG(wordID, chatID, resp.Result.MessageID)
		err = s.tgClient.Edit(text, chatID, resp.Result.MessageID, true, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Edit(msg string, chatID, msgID int, withFormat bool, key *model.InlineKeyboardMarkup) error {
	return s.tgClient.Edit(msg, chatID, msgID, withFormat, key)
}
