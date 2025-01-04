package telegram

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func (t *Telegram) addNotificationTimeTG(ctx context.Context, msg *model.Message, arg string) string {
	arg = strings.TrimSpace(arg)

	timeToAdd, err := parseTime(arg)
	if err != nil {
		return model.AddTimeEmptyMSG
	}

	err = t.srv.AddNotificationTime(ctx, msg.From.Username, timeToAdd)
	if err != nil {
		logger.Error(ctx, "can't create notification time", "err", err, "time", arg, "user", msg.From.Username)
		if errors.Is(err, model.ErrAlreadyExists) {
			return model.AddTimeAlreadyExistsMSG
		}
		if errors.Is(err, model.ErrNotificationLimit) {
			return model.AddTimeLimitMSG
		}
		return model.AddTimeErrorMSG
	}

	logger.Info(ctx, "notification time created", "time", arg, "user", msg.From.Username)
	return fmt.Sprintf(model.AddTimeSuccessMSG, arg)
}

func parseTime(text string) (time.Time, error) {
	sp := strings.Split(text, ":")
	if len(sp) != 2 {
		return time.Time{}, errors.New("invalid time format")
	}
	hours, err := strconv.Atoi(sp[0])
	if err != nil || hours < 0 || hours > 23 {
		return time.Time{}, errors.New("invalid time format")
	}

	minutes, err := strconv.Atoi(sp[1])
	if err != nil || minutes < 0 || minutes > 59 {
		return time.Time{}, errors.New("invalid time format")
	}

	res, _ := time.Parse(time.TimeOnly, "00:00:00")
	res = res.Add(time.Duration(hours)*time.Hour + time.Duration(minutes)*time.Minute)
	return res, nil
}
