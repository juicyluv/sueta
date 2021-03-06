package main

import (
	"context"
	"errors"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/juicyluv/sueta/post_service/app/config"
	"github.com/juicyluv/sueta/post_service/app/internal/server"
	"github.com/juicyluv/sueta/post_service/app/pkg/logger"
	"github.com/julienschmidt/httprouter"
)

var (
	configPath = flag.String("config-path", "app/config/config.yml", "path for application configuration file")
)

// @title SUETA User Service API
// @version 1.0.0
// @description API documentation for Sueta User Service. Navedi sueti, brat.

// @host localhost:8080
// @BasePath /api

func main() {
	flag.Parse()
	logger.Init()

	logger := logger.GetLogger()
	logger.Info("logger initialized")

	cfg := config.Get(*configPath, ".env")
	logger.Info("loaded config file")

	router := httprouter.New()
	logger.Info("initialized httprouter")

	logger.Info("connecting to database")

	logger.Info("starting the server")
	srv := server.NewServer(cfg, router, &logger)

	quit := make(chan os.Signal, 1)
	signals := []os.Signal{syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM}
	signal.Notify(quit, signals...)

	go func() {
		if err := srv.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf("cannot run the server: %v", err)
		}
	}()
	logger.Infof("server has been started on port %s", cfg.Http.Port)

	<-quit
	logger.Warn("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// TODO: disconnect to postgres
		logger.Info("closed mongo database connection")
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Errorf("server shutdown failed: %v", err)
	}

	logger.Info("server has been shutted down")
}
