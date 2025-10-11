package telegram

import (
	"context"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) disableReverseTG(ctx context.Context, username string) string {
	err := t.srv.SetReverseEnabled(ctx, username, false)
	if err != nil {
		logger.Error(ctx, "can't disable reverse language:", "user", username, "err", err)
		return model.DisableReverseErrorMSG
	}

	logger.Info(ctx, "reverse language disabled:", "user", username)
	return model.DisableReverseSuccessMSG
}
