package repository

import (
	"context"
	"fmt"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (r *Repository) CreateWord(ctx context.Context, word *model.Word) error {
	query := `
	INSERT INTO words (word, translated_word, example, translated_example,  username) 
	VALUES ($1, $2, $3, $4, $5)`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	_, err := r.conn.ExecContext(
		ctx, query,
		word.Word,
		word.TranslatedWord,
		word.Example,
		word.TranslatedExample,
		word.Username,
	)
	if err != nil {
		if IsPostgresError(err, model.PostgresUniqueConstraint) {
			return fmt.Errorf("CreateWord: %w", model.ErrAlreadyExists)
		}
		return fmt.Errorf("CreateWord: %w", err)
	}

	return nil
}
