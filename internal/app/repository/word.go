package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
)

func (r *Repository) CreateWord(ctx context.Context, word *model.Word) error {
	query := `
	INSERT INTO words (word, translated_word, example, translated_example, username) 
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

func (r *Repository) GetWordById(ctx context.Context, id int) (*model.Word, error) {
	query := `
	SELECT id, word, translated_word, example, translated_example, last_correct_answer
	FROM words WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	word := &model.Word{}
	err := r.conn.GetContext(ctx, word, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("GetWordById: %w", err)
	}

	return word, nil
}

func (r *Repository) GetWordsForNotification(ctx context.Context, username string, limit int) ([]*model.Word, error) {
	// TODO: подумать насчет коефа
	query := `
	SELECT id, word, translated_word, example, translated_example, last_correct_answer,
	       (correct_answer_count - COALESCE(CURRENT_DATE - last_correct_answer::date, 0)) AS koef
	FROM words
	WHERE username = $1
	ORDER BY koef
	LIMIT $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	words := make([]*model.Word, 0)
	err := r.conn.SelectContext(ctx, &words, query, username, limit)
	if err != nil {
		return nil, fmt.Errorf("GetWordsForNotification: %w", err)
	}

	return words, nil
}

func (r *Repository) AddCorrectAnswerToWord(ctx context.Context, id int) error {
	query := `
	UPDATE words
	SET correct_answer_count = correct_answer_count + 1,
	    last_correct_answer = NOW()
	WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	_, err := r.conn.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("AddCorrectAnswerToWord: %w", err)
	}

	return nil
}
