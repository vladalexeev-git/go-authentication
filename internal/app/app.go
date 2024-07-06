package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-authentication/config"
	v1 "go-authentication/internal/api/http/v1"
	"go-authentication/internal/repository"
	"go-authentication/internal/service"
	"go-authentication/pkg/JWT"
	"go-authentication/pkg/httpserver"
	"go-authentication/pkg/logger"
	"go-authentication/pkg/mongodb"
	"go-authentication/pkg/postgres"
	"log/slog"
	"os"
	"os/signal"
	//"sso/pkg/postgres"
	"go-authentication/pkg/utils"
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

	//Postgres
	l.Debug("postgres url from config", slog.String("postgres_url", cfg.Postgres.URL))

	pg, err := postgres.New(cfg.Postgres.URL, postgres.MaxPoolSize(cfg.Postgres.PoolMax))
	if err != nil {
		l.Error("can't connect to postgres", slog.String(utils.Operation, op), slog.String("error", err.Error()))
		return
	}
	defer pg.Close()

	//MongoDB
	mCl, err := mongodb.NewClient(cfg.MongoDB.URI, cfg.MongoDB.Username, cfg.MongoDB.Password)
	if err != nil {
		l.Error("can't connect to mongodb",
			slog.String("error", err.Error()))
		return
	}
	mDB := mCl.Database(cfg.MongoDB.DbName)

	// Repositories
	accountRepo := repository.NewAccountRepo(log, pg)
	sessionRepo := repository.NewSessionRepo(mDB, log)

	// Services
	accountService := service.NewAccountService(cfg, log, accountRepo, sessionRepo)
	sessionService := service.NewSessionService(cfg, sessionRepo)

	jwt, err := JWT.New(cfg.AccessToken.SigningKey, cfg.AccessToken.TTL)
	if err != nil {
		l.Error("can't create jwt token", slog.String("error", err.Error()))
		return
	}
	authService := service.NewAuthService(log, jwt, accountService, sessionService)

	// Handlers v1
	handler := gin.New()
	v1.SetupHandlers(handler, log, cfg, accountService, sessionService, authService)

	// HTTP Server
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	l.Info("server is up and running",
		slog.String("port", cfg.HTTP.Port))

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	select {
	case s := <-interrupt:
		l.Info("got signal",
			slog.String("signal", s.String()))
	case err := <-httpServer.Notify():
		l.Error("http server got error, shutting down...",
			slog.String("error", err.Error()))
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Error(fmt.Sprintf("app - Run - httpServer.Shutdown: %s", err.Error()))
	}
}
