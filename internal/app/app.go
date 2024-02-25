package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"sso/config"
	v1 "sso/internal/api/http/v1"
	"sso/pkg/httpserver"
	"sso/pkg/logger"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	log := logger.SetupLogger(cfg.Logger)

	// HTTP Server
	handler := gin.New()

	log.Info("listen and serve...", slog.String("port", cfg.HTTP.Port),
		slog.String("log level", cfg.Logger.LogLevel))

	v1.SetupHandlers(handler, log, cfg)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		//l.Info("app - Run - signal: " + s.String())
		log.Info(fmt.Sprintf("app - Run - signal: %s", s.String()))
	case err := <-httpServer.Notify():
		log.Error(fmt.Sprintf("app - Run - httpServer.Notify: %s", err.Error()))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err.Error()))
	}
}
