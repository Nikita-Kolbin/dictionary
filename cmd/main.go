package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Nikita-Kolbin/dictionary/internal/app/config"
	"github.com/Nikita-Kolbin/dictionary/internal/app/repository"
	"github.com/Nikita-Kolbin/dictionary/internal/app/service"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/clients/telegram"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
)

func main() {
	ctx := context.Background()
	if err := initApp(ctx); err != nil {
		logger.Error(ctx, "failed to init app", "error", err)
	}
}

func initApp(ctx context.Context) error {
	cfg, err := config.New()
	if err != nil {
		return fmt.Errorf("init config failed: %w", err)
	}

	repo, err := repository.New(ctx, &cfg.Postgres)
	if err != nil {
		return fmt.Errorf("init reposytory failed: %w", err)
	}
	defer repo.Close(ctx)

	tgCli := telegram.New(cfg.TelegramToken)

	srv := service.New(repo, tgCli)

	srv.RunTelegramProcessor(ctx)

	ch := make(chan os.Signal, 1)
	<-ch

	return fmt.Errorf("failed to init app")
}