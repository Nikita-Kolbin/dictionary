package telegram

import (
	"context"
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
	return t.srv.BuildWordMessage(word), word.ID
}
