package httpserver

import (
	"context"
	"net/http"

	"github.com/Nikita-Kolbin/dictionary/internal/app/config"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/logger"
	"github.com/Nikita-Kolbin/dictionary/internal/pkg/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Server struct {
	server *http.Server
}

func New(cfg *config.ListenerConfig) *Server {
	server := http.NewServeMux()

	// Create non-global registry.
	reg := prometheus.NewRegistry()

	// Add go runtime metrics and process collectors.
	reg.MustRegister(
		metrics.New()...,
	)

	// Expose /metrics HTTP endpoint using the created custom registry.
	server.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	srv := &http.Server{
		Addr:         cfg.GetHostPort(),
		Handler:      server,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return &Server{
		server: srv,
	}
}

func (s *Server) Run(ctx context.Context) {
	go func(ctx context.Context) {
		if err := s.server.ListenAndServe(); err != nil {
			logger.Fatal(ctx, "run server error", "error", err)
		}
	}(ctx)
}
