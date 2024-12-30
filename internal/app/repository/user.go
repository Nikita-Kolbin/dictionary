package repository

import (
	"context"
	"fmt"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
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
