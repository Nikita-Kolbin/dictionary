package telegram

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) getOneWordTG(ctx context.Context, username string) (string, int) {
	word, err := t.srv.GetOneWord(ctx, username)
	if err != nil {
		logger.Error(ctx, "can't get one word:", "user", username, "err", err)
		return model.GetErrorMSG, 0
	}
	if word == nil {
		return model.GetUserHaveNotWordsMSG, 0
	}
	logger.Info(ctx, "word given:", "user", username, "word", word.Word)
	return buildWordMessage(word), word.ID
}

func buildWordMessage(word *model.Word) string {
	if word == nil {
		return ""
	}

	builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf(model.GetSuccessWordMSG, escapeFormatChars(word.Word)))
	builder.WriteRune('\n')
	builder.WriteString(fmt.Sprintf(model.GetSuccessTranslateMSG, escapeFormatChars(word.TranslatedWord)))
	if len(word.Example) > 0 {
		builder.WriteRune('\n')
		builder.WriteString(fmt.Sprintf(model.GetSuccessExampleMSG, escapeFormatChars(word.Example)))
	}
	if len(word.TranslatedExample) > 0 {
		builder.WriteRune('\n')
		builder.WriteString(fmt.Sprintf(model.GetSuccessExampleTranslateMSG, escapeFormatChars(word.TranslatedExample)))
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
