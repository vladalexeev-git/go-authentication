package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"os"
	"os/signal"
	"sso/config"
	v1 "sso/internal/api/http/v1"
	"sso/internal/repository"
	"sso/internal/service"
	"sso/pkg/httpserver"
	"sso/pkg/logger"
	"sso/pkg/postgres"
	"sso/pkg/utils"
	"syscall"
)

// Run creates objects via constructors.
func Run(cfg *config.Config) {
	const op = "app.run"

	log := logger.SetupLogger(cfg.Logger)
	log.Info("app is starting...",
		slog.String("log_environment", cfg.Logger.Env))

	// Postgres
	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	log.Debug("postgres url from config", slog.String("postgres_url", cfg.Postgres.URL))
	//defer pg.Close()
	if err != nil {
		log.Error("can't connect to postgres", slog.String(utils.Operation, op), slog.String("error", err.Error()))
		return
	}

	accountRepo := repository.NewAccountRepo(log, pg)

	// Services
	accountService := service.NewAccountService(cfg, log, accountRepo)

	//Handlers v1
	handler := gin.New()
	v1.SetupHandlers(handler, log, cfg, accountService)

	// HTTP Server
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	log.Info("server is up and running",
		slog.String("port", cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		//l.Info("app - Run - signal: " + s.String())
		log.Info("got signal",
			slog.String(utils.Operation, op),
			slog.String("signal", s.String()))
	case err := <-httpServer.Notify():
		log.Error("http server got error, shutting down...",
			slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		log.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err.Error()))
	}
}
