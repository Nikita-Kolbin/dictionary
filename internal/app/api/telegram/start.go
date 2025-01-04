package telegram

import (
	"context"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) createUserTG(ctx context.Context, msg *model.Message) string {
	user := &model.User{
		Username: msg.From.Username,
		ChatID:   msg.Chat.ID,
	}
	if err := t.srv.CreateUser(ctx, user); err != nil {
		logger.Error(ctx, "can't, create user", "err", err, "user", user.Username)
		return model.StartMSG
	}
	logger.Info(ctx, "user created", "user", user.Username)
	return model.StartMSG
}
