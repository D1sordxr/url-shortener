package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/D1sordxr/url-shortener/internal/application/url/usecase"
	loadApp "github.com/D1sordxr/url-shortener/internal/infrastructure/app"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/config"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/logger"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/url/repo"
	"github.com/D1sordxr/url-shortener/internal/transport/http"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler"
	"github.com/D1sordxr/url-shortener/internal/transport/http/middleware"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.NewApiConfig()

	log := logger.New(defaultLogger)
	log.Debug("Config data", "config", cfg)

	storageConn, err := dbpg.New(cfg.Storage.ConnectionString(), nil, nil)
	if err != nil {
		log.Error("Failed to connect to database", "error", err.Error())
		return
	}
	defer func() { _ = storageConn.Master.Close() }()
	if err = postgres.SetupStorage(storageConn.Master, cfg.Storage); err != nil {
		log.Error("Failed to setup storage", "error", err.Error())
		return
	}

	urlRepository := repo.NewRepository(log, storageConn)
	urlUC := usecase.New(log, urlRepository)
	urlHandler := handler.New(urlUC)
	urlRouteRegisterer := url.NewRouteRegisterer(urlHandler, middleware.Stat)

	httpServer := http.NewServer(
		log,
		&cfg.Server,
		urlRouteRegisterer,
	)

	app := loadApp.NewApp(
		log,
		httpServer,
	)
	app.Run(ctx)
}

var defaultLogger zerolog.Logger

func init() {
	zlog.Init()
	defaultLogger = zlog.Logger
}
