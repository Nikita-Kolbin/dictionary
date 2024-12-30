package model

import "time"

type Word struct {
	ID                 int        `db:"id"`
	Word               string     `db:"word"`
	TranslatedWord     string     `db:"translated_word"`
	Example            string     `db:"example"`
	TranslatedExample  string     `db:"translated_example"`
	Username           string     `db:"username"`
	CorrectAnswerCount int        `db:"correct_answer_count"`
	LastCorrectAnswer  *time.Time `db:"last_correct_answer"`
	Created            time.Time  `db:"created"`
	Koef               int        `db:"koef"`
}
