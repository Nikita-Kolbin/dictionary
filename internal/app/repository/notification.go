package repository

import (
	"context"
	"fmt"
	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"time"
)

func (r *Repository) GetNotificationTimes(ctx context.Context, username string) ([]time.Time, error) {
	query := `SELECT time FROM notification_times WHERE username = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	strTimes := make([]string, 0)

	err := r.conn.SelectContext(ctx, &strTimes, query, username)
	if err != nil {
		return nil, fmt.Errorf("GetNotificationTimes: %w", err)
	}

	times := make([]time.Time, 0, len(strTimes))
	for _, strTime := range strTimes {
		t, _ := time.Parse(time.TimeOnly, strTime)
		times = append(times, t)
	}

	return times, nil
}

func (r *Repository) AddNotificationTime(ctx context.Context, username string, t time.Time) error {
	query := `INSERT INTO notification_times (username, time) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	_, err := r.conn.ExecContext(ctx, query, username, t)
	if err != nil {
		if IsPostgresError(err, model.PostgresUniqueConstraint) {
			return fmt.Errorf("AddNotificationTimes: %w", model.ErrAlreadyExists)
		}
		return fmt.Errorf("AddNotificationTimes: %w", err)
	}

	return nil
}
