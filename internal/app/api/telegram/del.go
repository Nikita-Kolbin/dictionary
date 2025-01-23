package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) delWordTG(ctx context.Context, msg *model.Message, arg string) string {
	arg = strings.TrimSpace(arg)

	err := t.srv.DeleteWord(ctx, arg, msg.From.Username)
	if err != nil {
		logger.Error(ctx, "can't delete word", "err", err, "word", arg, "user", msg.From.Username)
		if errors.Is(err, model.ErrNotFound) {
			return model.DelEmptyMSG
		}
		return model.DelErrorMSG
	}

	logger.Info(ctx, "word deleted", "word", arg, "user", msg.From.Username)
	return fmt.Sprintf(model.DelSuccessMSG, arg)
}
