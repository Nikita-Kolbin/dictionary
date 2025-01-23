package telegram

import (
	"context"
	"errors"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) backupTG(ctx context.Context, msg *model.Message) (filePath, text string) {
	filePath, err := t.srv.MakeUserWordsBackup(ctx, msg.From.Username)
	if err != nil {
		logger.Error(ctx, "can't, make backup", "err", err)
		if errors.Is(err, model.ErrBackupLimit) {
			return "", model.BackupLimit
		}
		return "", model.BackupErrorMSG
	}
	return filePath, ""
}
