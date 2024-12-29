package main

import (
	"context"
	"fmt"
	"github.com/Nikita-Kolbin/dictionary/internal/app/config"
	"github.com/Nikita-Kolbin/dictionary/internal/app/service"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/clients/telegram"
	"os"

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

	tgCli := telegram.New(cfg.TelegramToken)

	srv := service.New(tgCli)

	srv.RunTelegramProcessor(ctx)

	ch := make(chan os.Signal, 1)
	<-ch

	return fmt.Errorf("failed to init app")
}
