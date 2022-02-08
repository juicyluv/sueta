package server

import (
	"context"
	"net/http"
	"time"

	"github.com/juicyluv/sueta/user_service/app/config"
	"github.com/juicyluv/sueta/user_service/app/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

// Server represents http server.
type Server struct {
	server *http.Server
	logger *logger.Logger
}

// NewServer returns a new Server instance.
func NewServer(cfg *config.Config, handler *httprouter.Router, logger *logger.Logger) *Server {
	return &Server{
		server: &http.Server{
			Handler:        handler,
			WriteTimeout:   time.Duration(cfg.Http.WriteTimeout) * time.Second,
			ReadTimeout:    time.Duration(cfg.Http.ReadTimeout) * time.Second,
			MaxHeaderBytes: cfg.Http.MaxHeaderBytes << 20,
			Addr:           ":" + cfg.Http.Port,
		},
		logger: logger,
	}
}

// Run starts http server. Returns an error on failure.
func (s *Server) Run() error {
	return s.server.ListenAndServe()
}

// Shutdown closes all connections and shuts down http server.
// It uses httpServer.Shutdown() function. Returns an error on failure.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
