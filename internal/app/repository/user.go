package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/lib/pq"
)

func (r *Repository) CreateUser(ctx context.Context, user *model.User) error {
	query := `INSERT INTO users (username, chat_id) VALUES ($1, $2)`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	_, err := r.conn.ExecContext(ctx, query, user.Username, user.ChatID)
	if err != nil {
		if IsPostgresError(err, model.PostgresUniqueConstraint) {
			return model.ErrAlreadyExists
		}
		return fmt.Errorf("CreateUser: %w", err)
	}

	return nil
}

func (r *Repository) GetUser(ctx context.Context, username string) (*model.User, error) {
	query := `SELECT username, chat_id, notification_word_count, created FROM users WHERE username = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	user := &model.User{}

	err := r.conn.GetContext(ctx, user, query, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("GetUser: %w", err)
	}

	return user, nil
}

func (r *Repository) GetUsers(ctx context.Context, usernames []string) ([]*model.User, error) {
	query := `SELECT username, chat_id, notification_word_count, created FROM users WHERE username = ANY($1)`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	users := make([]*model.User, 0)

	err := r.conn.SelectContext(ctx, &users, query, pq.Array(usernames))
	if err != nil {
		return nil, fmt.Errorf("GetUsers: %w", err)
	}

	return users, nil
}

func (r *Repository) SetWordsCount(ctx context.Context, username string, count int) error {
	query := `UPDATE users SET notification_word_count = $1 WHERE username = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	res, err := r.conn.ExecContext(ctx, query, count, username)
	if err != nil {
		return fmt.Errorf("SetWordsCount: %w", err)
	}
	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return model.ErrNotFound
	}

	return nil
}
