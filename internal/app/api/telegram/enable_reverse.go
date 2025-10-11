package telegram

import (
	"context"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) enableReverseTG(ctx context.Context, username string) string {
	err := t.srv.SetReverseEnabled(ctx, username, true)
	if err != nil {
		logger.Error(ctx, "can't enable reverse language:", "user", username, "err", err)
		return model.EnableReverseErrorMSG
	}

	logger.Info(ctx, "reverse language enabled:", "user", username)
	return model.EnableReverseSuccessMSG
}
