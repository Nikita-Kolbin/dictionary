package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) delNotificationTimeTG(ctx context.Context, msg *model.Message, arg string) string {
	arg = strings.TrimSpace(arg)

	timeToDelete, err := parseTime(arg)
	if err != nil {
		return model.DelTimeEmptyMSG
	}

	err = t.srv.DelNotificationTime(ctx, msg.From.Username, timeToDelete)
	if err != nil {
		logger.Error(ctx, "can't delete notification time", "err", err, "time", arg, "user", msg.From.Username)
		if errors.Is(err, model.ErrNotFound) {
			return model.DelTimeEmptyMSG
		}
		return model.DelTimeErrorMSG
	}

	logger.Info(ctx, "notification time deleted", "time", arg, "user", msg.From.Username)
	return fmt.Sprintf(model.DelTimeSuccessMSG, arg)
}
