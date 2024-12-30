package service

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
	
	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func parseWord(text string) (word, trWord, example, trExample string) {
	// TODO: Добавить разделителей
	sp := strings.Split(text, ",")

	if len(sp) > 0 {
		word = strings.TrimSpace(sp[0])
	}
	if len(sp) > 1 {
		trWord = strings.TrimSpace(sp[1])
	}
	if len(sp) > 2 {
		example = strings.TrimSpace(sp[2])
	}
	if len(sp) > 3 {
		trExample = strings.TrimSpace(sp[3])
	}

	return
}

func parseTime(text string) (time.Time, error) {
	sp := strings.Split(text, ":")
	if len(sp) != 2 {
		return time.Time{}, errors.New("invalid time format")
	}
	hours, err := strconv.Atoi(sp[0])
	if err != nil || hours < 0 || hours > 23 {
		return time.Time{}, errors.New("invalid time format")
	}

	minutes, err := strconv.Atoi(sp[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return time.Time{}, errors.New("invalid time format")
	}

	res, _ := time.Parse(time.TimeOnly, "00:00:00")
	res = res.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)
	return res, nil
}

func (s *Service) getKeyTG(wordID, chatID, msgID int) *model.InlineKeyboardMarkup {
	goodData := &model.CallbackData{
		WordID:    wordID,
		ChatID:    chatID,
		MessageID: msgID,
		Correct:   true,
	}
	badData := &model.CallbackData{
		WordID:    wordID,
		ChatID:    chatID,
		MessageID: msgID,
		Correct:   false,
	}

	good, _ := json.Marshal(goodData)
	bad, _ := json.Marshal(badData)

	key := &model.InlineKeyboardMarkup{
		InlineKeyboard: [][]*model.InlineKeyboardButton{{
			{
				Text:         model.GoodButton,
				CallbackData: string(good),
			},
			{
				Text:         model.BadButton,
				CallbackData: string(bad),
			},
		}},
	}

	return key
}
