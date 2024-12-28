package httpserver

import (
	"net/http"
	"time"
)

type Server struct {
	server *http.Server
}

func New(address string, handler http.Handler, readTimeout, writeTimeout, idleTimeout time.Duration) *Server {
	srv := &http.Server{
		Addr:         address,
		Handler:      handler,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
		IdleTimeout:  idleTimeout,
	}

	return &Server{
		server: srv,
	}
}

func (s *Server) Run() error {
	return s.server.ListenAndServe()
}
