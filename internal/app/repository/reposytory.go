package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/Nikita-Kolbin/dictionary/internal/app/config"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
	_ "github.com/jackc/pgx/v5/stdlib" // driver
	"github.com/jmoiron/sqlx"
)

const (
	timeBetweenPings = 5 * time.Second
	queryTimeout     = 5 * time.Second
)

type Repository struct {
	conn    *sqlx.DB
	timeout time.Duration
}

func New(ctx context.Context, cfg *config.PostgresConfig) (*Repository, error) {
	var err error
	repo := &Repository{
		timeout: cfg.DBTimeout,
	}

	logger.Info(ctx, "creating postgres connection")
	repo.conn, err = sqlx.Open("pgx", cfg.GetDSN())
	if err != nil {
		return nil, fmt.Errorf("failed to connect with postgres credentials")
	}
	repo.conn.SetMaxOpenConns(cfg.DBMaxConn)
	repo.conn.SetConnMaxLifetime(cfg.DBMaxConnLifeTime)
	repo.conn.SetConnMaxIdleTime(cfg.DBMaxConnIdleTime)
	if !pingPG(repo.conn) {
		return nil, fmt.Errorf("failed to ping postgres credentials: %w", err)
	}

	return repo, nil
}

func pingPG(conn *sqlx.DB) bool {
	var err error
	for i := 0; i < 5; i++ {
		if err = conn.Ping(); err == nil {
			return true
		}
		logger.Error(context.Background(), "failed to ping postgres", err)
		time.Sleep(timeBetweenPings)
	}
	return false
}

func (r *Repository) Close(ctx context.Context) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	logger.Info(ctx, "closing master postgres connection pool")
	_ = r.conn.Close()
}
