package main

import (
	"context"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/logger"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"

	loadApp "github.com/D1sordxr/url-shortener/internal/infrastructure/app"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/config"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/worker"

	rabbitAdapter "github.com/D1sordxr/url-shortener/internal/infrastructure/broker/rabbitmq/notification"
	"github.com/D1sordxr/url-shortener/internal/infrastructure/cache/redis/notification"
	notificationRepository "github.com/D1sordxr/url-shortener/internal/infrastructure/storage/postgres/repositories/notification"
	workerHandler "github.com/D1sordxr/url-shortener/internal/transport/rabbitmq/notification/handler"

	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/zlog"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	cfg := config.NewWorkerConfig()

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
	notificationCacheAdapter := notification.NewAdapter(cacheConn)

	brokerConn, err := rabbitmq.Connect(cfg.Broker.GetConnectionString(), 3, time.Second*3)
	if err != nil {
		log.Error("Broker connection failure", "error", err.Error())
	}
	defer func() { _ = brokerConn.Close() }()

	notificationQueue, err := rabbitAdapter.NewQueue(log, brokerConn, cfg.Broker.DeclareExchange)
	if err != nil {
		log.Error("RabbitMQ queue failure", "error", err.Error())
		return
	}
	notificationProducer := rabbitAdapter.NewProducer(notificationQueue)
	notificationConsumer := rabbitAdapter.NewConsumer(log, notificationQueue)

	notificationWriter := workerHandler.NewNotificationWriter(
		log,
		notificationRepo,
		notificationProducer,
		notificationCacheAdapter,
	)
	notificationProcessor := workerHandler.NewProcessor(log, notificationConsumer)
	notificationWorker := worker.NewWorker(
		log,
		notificationProcessor,
		notificationWriter,
	)

	app := loadApp.NewApp(
		log,
		notificationWorker,
	)
	app.Run(ctx)
}

var defaultLogger zerolog.Logger

func init() {
	zlog.Init()
	defaultLogger = zlog.Logger
}
