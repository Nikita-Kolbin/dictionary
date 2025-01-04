package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) setWordCountTG(ctx context.Context, msg *model.Message, arg string) string {
	arg = strings.TrimSpace(arg)

	cnt, err := strconv.Atoi(arg)
	if err != nil {
		return model.SetCountEmptyMSG
	}
	if cnt <= 0 || cnt > 25 {
		return model.SetCountEmptyMSG
	}

	err = t.srv.SetWordsCount(ctx, msg.From.Username, cnt)
	if err != nil {
		logger.Error(ctx, "can't set word count", "err", err, "count", cnt, "user", msg.From.Username)
		if errors.Is(err, model.ErrNotFound) {
			return model.SetCountUserNotFoundMSG
		}
		return model.SetCountErrorMSG
	}

	logger.Info(ctx, "word count set", "count", cnt, "user", msg.From.Username)
	return fmt.Sprintf(model.SetCountSuccessMSG, cnt)
}
