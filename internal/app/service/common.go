package service

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func parseWord(text string) (word, trWord, example, trExample string) {
	// TODO: Добавить разделителей
	sp := strings.Split(text, ",")

	if len(sp) > 0 {
		word = strings.TrimSpace(sp[0])
	}
	if len(sp) > 1 {
		trWord = strings.TrimSpace(sp[1])
	}
	if len(sp) > 2 {
		example = strings.TrimSpace(sp[2])
	}
	if len(sp) > 3 {
		trExample = strings.TrimSpace(sp[3])
	}

	return
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
