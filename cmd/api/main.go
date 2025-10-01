package main

import (
	"context"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/logger"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/url"
	handler2 "github.com/D1sordxr/url-shortener/internal/transport/http/api/url/handler"
	"github.com/D1sordxr/url-shortener/internal/transport/http/middleware"
	"os"
	"os/signal"
	"syscall"

	loadApp "github.com/D1sordxr/url-shortener/internal/infrastructure/app"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/config"
	"github.com/D1sordxr/url-shortener/internal/transport/http"

	notificationUseCase "github.com/D1sordxr/url-shortener/internal/application/notification/usecase"
	notificationCache "github.com/D1sordxr/url-shortener/internal/infrastructure/cache/redis/notification"
	notificationRepository "github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/notification"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/notify"
	"github.com/D1sordxr/url-shortener/internal/transport/http/api/notify/handler"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/redis"
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
	notificationRepo := notificationRepository.NewRepository(log, storageConn)

	cacheConn := redis.New(cfg.Cache.ClientAddress, cfg.Cache.Password, 1)
	defer func() { _ = cacheConn.Close() }()
	notificationCacheAdapter := notificationCache.NewAdapter(cacheConn)

	notificationUC := notificationUseCase.NewUseCase(
		log,
		notificationCacheAdapter,
		notificationRepo,
	)

	notificationHandlers := handler.NewHandlers(notificationUC)
	notificationRouteRegisterer := notify.NewRouteRegisterer(notificationHandlers)

	urlHandler := handler2.NewHandler()
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
