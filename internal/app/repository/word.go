package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
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

func (r *Repository) GetWordByID(ctx context.Context, id int) (*model.Word, error) {
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
		return nil, fmt.Errorf("GetWordByID: %w", err)
	}

	return word, nil
}

func (r *Repository) DeleteWord(ctx context.Context, word, username string) error {
	query := `DELETE FROM words WHERE word = $1 AND username = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	res, err := r.conn.ExecContext(ctx, query, word, username)
	if err != nil {
		return fmt.Errorf("DeleteWord: %w", err)
	}
	if cnt, _ := res.RowsAffected(); cnt == 0 {
		return fmt.Errorf("DeleteWord: %w", model.ErrNotFound)
	}

	return nil
}

func (r *Repository) GetWordsForNotification(ctx context.Context, username string, limit int) ([]*model.Word, error) {
	// TODO: подумать насчет коефа
	query := `
	SELECT id, word, translated_word, example, translated_example, last_correct_answer,
		(correct_answer_count - COALESCE(CURRENT_DATE - last_correct_answer::date, 0)) AS koef,
		((SELECT reverse_enabled FROM users WHERE username = $1) AND last_answer_is_original) AS need_reverse_lang
	FROM words
	WHERE username = $1 AND current_message_id IS NULL
	ORDER BY koef, RANDOM()
	LIMIT $2`

	querySetLastLang := `
	UPDATE words 
	SET last_answer_is_original = $1, last_send_date = now() 
	WHERE id = any($2)`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	words := make([]*model.Word, 0)
	err := r.conn.SelectContext(ctx, &words, query, username, limit)
	if err != nil {
		return nil, fmt.Errorf("GetWordsForNotification: %w", err)
	}

	reverseLangList := make([]int, 0)
	originLangList := make([]int, 0)
	for _, word := range words {
		if word.NeedReverseLang {
			reverseLangList = append(reverseLangList, word.ID)
		} else {
			originLangList = append(originLangList, word.ID)
		}
	}

	// Русские слова
	if len(reverseLangList) != 0 {
		_, err = r.conn.ExecContext(ctx, querySetLastLang, false, reverseLangList)
		if err != nil {
			logger.Error(ctx, "GetWordsForNotification: can't change last answer language", err)
		}
	}

	// Английские слова
	if len(originLangList) != 0 {
		_, err = r.conn.ExecContext(ctx, querySetLastLang, true, originLangList)
		if err != nil {
			logger.Error(ctx, "GetWordsForNotification: can't change last answer language", err)
		}
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

func (r *Repository) GetAllUserWords(ctx context.Context, username string) ([]*model.Word, error) {
	query := `
	SELECT id, word, translated_word, example, translated_example, 
	       correct_answer_count, last_correct_answer, created
	FROM words
	WHERE username = $1
	ORDER BY id`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	words := make([]*model.Word, 0)
	err := r.conn.SelectContext(ctx, &words, query, username)
	if err != nil {
		return nil, fmt.Errorf("GetAllUserWords: %w", err)
	}

	return words, nil
}

func (r *Repository) UpdateWordCurrentMessageID(ctx context.Context, wordID int, msgID *int) error {
	query := `UPDATE words SET current_message_id = $1 WHERE id = $2`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	_, err := r.conn.ExecContext(ctx, query, msgID, wordID)
	if err != nil {
		return fmt.Errorf("UpdateWordCurrentMessageID: %w", err)
	}

	return nil
}

func (r *Repository) GetOldSendWords(ctx context.Context) ([]*model.Word, error) {
	const wordTTL = time.Hour * 24

	query := `
	SELECT w.id, w.word, w.translated_word, w.example, 
	       w.translated_example, w.correct_answer_count, 
	       w.last_correct_answer, w.created, 
	       w.current_message_id, u.chat_id
	FROM words as w
	JOIN users as u ON u.username = w.username
	WHERE current_message_id IS NOT NULL AND now() - w.last_send_date > $1`

	ctx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	words := make([]*model.Word, 0)
	err := r.conn.SelectContext(ctx, &words, query, wordTTL)
	if err != nil {
		return nil, fmt.Errorf("GetOldOpenedWords: %w", err)
	}

	return words, nil
}
