package service

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/model"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
	"github.com/xuri/excelize/v2"
)

func (s *Service) MakeUserWordsBackup(ctx context.Context, username string) (path string, err error) { //nolint
	user, err := s.repo.GetUser(ctx, username)
	if err != nil {
		return "", err
	}
	if user.LastBackup.After(time.Now().Add(-24 * time.Hour)) {
		return "", model.ErrBackupLimit
	}

	words, err := s.repo.GetAllUserWords(ctx, username)
	if err != nil {
		return "", err
	}

	file := excelize.NewFile()
	_ = file.SetSheetName("Sheet1", "Words")

	// headers
	_ = file.SetCellValue("Words", "A1", "ID")
	_ = file.SetCellValue("Words", "B1", "Слово")
	_ = file.SetCellValue("Words", "C1", "Перевод")
	_ = file.SetCellValue("Words", "D1", "Пример")
	_ = file.SetCellValue("Words", "E1", "Перевод примера")
	_ = file.SetCellValue("Words", "F1", "Кол-во ответов")
	_ = file.SetCellValue("Words", "G1", "Дата добавления")

	// filling
	for i, word := range words {
		row := i + 2
		_ = file.SetCellValue("Words", fmt.Sprintf("A%d", row), word.ID)
		_ = file.SetCellValue("Words", fmt.Sprintf("B%d", row), word.Word)
		_ = file.SetCellValue("Words", fmt.Sprintf("C%d", row), word.TranslatedWord)
		_ = file.SetCellValue("Words", fmt.Sprintf("D%d", row), word.Example)
		_ = file.SetCellValue("Words", fmt.Sprintf("E%d", row), word.TranslatedExample)
		_ = file.SetCellValue("Words", fmt.Sprintf("F%d", row), word.CorrectAnswerCount)
		_ = file.SetCellValue("Words", fmt.Sprintf("G%d", row), word.Created)
	}

	dir := "./temp"
	_ = os.MkdirAll(dir, os.ModePerm)
	path = fmt.Sprintf("%s/%s.xlsx", dir, username)

	err = file.SaveAs(path)
	if err != nil {
		return "", err
	}

	// remove file after 1 minute
	go func() {
		time.Sleep(1 * time.Minute)
		err = os.Remove(path)
		if err != nil {
			logger.Error(ctx, "remove temp excel file error", err, "path", path)
		}
	}()

	err = s.repo.UpdateUserLastBackup(ctx, username)
	if err != nil {
		logger.Error(ctx, "update user backup time error", err, "username", username)
	}

	return path, nil
}
