package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) addWordTG(ctx context.Context, msg *model.Message, arg string) string {
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

	err := t.srv.CreateWord(ctx, word)
	if err != nil {
		logger.Error(ctx, "can't, create word", "err", err, "word", w, "user", msg.From.Username)
		if errors.Is(err, model.ErrAlreadyExists) {
			return model.AddAlreadyExistsMSG
		}
		return model.AddErrorMSG
	}

	logger.Info(ctx, "word created", "word", w, "user", msg.From.Username)
	return fmt.Sprintf(model.AddSuccessMSG, w)
}

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
