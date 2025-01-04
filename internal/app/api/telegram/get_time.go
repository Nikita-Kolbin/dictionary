package telegram

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) getNotificationTimeTG(ctx context.Context, msg *model.Message) string {
	times, err := t.srv.GetNotificationTimes(ctx, msg.From.Username)
	if err != nil {
		logger.Error(ctx, "can't get notification time", "err", err, "user", msg.From.Username)
		return model.GetTimeErrorMSG
	}
	if len(times) == 0 {
		return model.GetTimeEmptyMSG
	}

	logger.Info(ctx, "notification time given", "user", msg.From.Username)
	return buildNotificationTimeMessage(times)
}

func buildNotificationTimeMessage(times []time.Time) string {
	builder := strings.Builder{}
	builder.WriteString(model.GetTimeSuccessMSG)
	for _, t := range times {
		builder.WriteRune('\n')
		strTime := fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
		builder.WriteString(strTime)
	}
	return builder.String()
}
