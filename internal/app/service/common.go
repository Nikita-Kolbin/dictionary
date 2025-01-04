package service

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func getKeyTG(wordID, chatID, msgID int) *model.InlineKeyboardMarkup {
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

func buildWordMessage(word *model.Word) string {
	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(model.GetSuccessWordMSG, word.Word))
	builder.WriteRune('\n')
	builder.WriteString(fmt.Sprintf(model.GetSuccessTranslateMSG, word.TranslatedWord))
	if len(word.Example) > 0 {
		builder.WriteRune('\n')
		builder.WriteString(fmt.Sprintf(model.GetSuccessExampleMSG, word.Example))
	}
	if len(word.TranslatedExample) > 0 {
		builder.WriteRune('\n')
		builder.WriteString(fmt.Sprintf(model.GetSuccessExampleTranslateMSG, word.TranslatedExample))
	}

	query := url.Values{}
	query.Add("sl", "eng")
	query.Add("tl", "rus")
	query.Add("text", word.Word)
	translatorURL := url.URL{
		Scheme:   "https",
		Host:     "www.reverso.net",
		Path:     "перевод-текста",
		Fragment: query.Encode(),
	}
	builder.WriteRune('\n')
	builder.WriteString(fmt.Sprintf(model.GetSuccessOpenInTranslator, translatorURL.String()))

	return builder.String()
}
