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

	// Logger
	log := logger.SetupLogger(cfg.Logger)
	l := log.With(slog.String(utils.Operation, op))
	l.Info("app is starting...",
		slog.String("log_environment", cfg.Logger.Env))

	// Postgres
	l.Debug("postgres url from config", slog.String("postgres_url", cfg.Postgres.URL))

	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		l.Error("can't connect to postgres", slog.String(utils.Operation, op), slog.String("error", err.Error()))
		return
	}
	defer pg.Close()

	// MongoDB
	//mCl, err := mongodb.NewClient(cfg.MongoDB.URI, cfg.MongoDB.Username, cfg.MongoDB.Password)
	//if err != nil {
	//	l.Error("can't connect to mongodb",
	//		slog.String("error", err.Error()))
	//	return
	//}
	//mDB := mCl.Database(cfg.MongoDB.Database)
	//
	//// Repositories
	accountRepo := repository.NewAccountRepo(l, pg)
	//sessionRepo := repository.NewSessionRepo(mDB)

	// Services
	accountService := service.NewAccountService(cfg, l, accountRepo)

	// Handlers v1
	handler := gin.New()
	v1.SetupHandlers(handler, l, cfg, accountService)

	// HTTP Server
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	l.Info("server is up and running",
		slog.String("port", cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		//l.Info("app - Run - signal: " + s.String())
		l.Info("got signal",
			//slog.String(utils.Operation, op),
			slog.String("signal", s.String()))
	case err := <-httpServer.Notify():
		l.Error("http server got error, shutting down...",
			//slog.String(utils.Operation, op),
			slog.String("error", err.Error()))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err.Error()))
	}
}
