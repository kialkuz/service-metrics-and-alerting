package app

import (
	"context"
	"fmt"
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/file"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/memory"
	migrations "kialkuz/service-metrics-and-alerting/internal/infrastructure/storage"
	dbWrapper "kialkuz/service-metrics-and-alerting/internal/infrastructure/storage/postgresql"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"kialkuz/service-metrics-and-alerting/internal/server/handler"
	pkgContracts "kialkuz/service-metrics-and-alerting/pkg/contracts"

	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/publisher"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	service "kialkuz/service-metrics-and-alerting/internal/service/server"
	auditlogs "kialkuz/service-metrics-and-alerting/internal/service/server/auditlogs"
	auditLogsQueue "kialkuz/service-metrics-and-alerting/internal/service/server/auditlogs/queue"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	Handler       *handler.MetricsHandler
	Pool          *pgxpool.Pool
	LogsSaver     *service.LogsSaverService
	LogsPublisher publisher.Publisher
}

func NewApp(ctx context.Context, logger *zap.SugaredLogger, cfg *appConfig.Config) (*App, error) {
	var (
		storage service.MetricsRepository
		pinger  *db.MemStorage
		pool    *pgxpool.Pool
	)

	fileStorage, err := file.NewFileStorage(cfg)
	if err != nil {
		return nil, fmt.Errorf("error init fileStorage: %s", err.Error())
	}

	if cfg.DB.DatabaseURI != "" {
		var err error

		pool, err = pgxpool.New(ctx, cfg.DB.DatabaseURI)
		if err != nil {
			return nil, err
		}

		dbStorage := db.NewDBStorage(dbWrapper.NewDB(pool), pool)

		storage = dbStorage
		pinger = dbStorage

		migrations.RunMigrations("pgx", cfg.DB.DatabaseURI)
	} else {
		storage = memory.NewMemoryStorage()
	}

	auditLogsQueue := auditLogsQueue.NewAuditLogsQueryList(map[string]int{
		subscriber.AuditTypeFile: auditLogsQueue.MaxThreads,
		subscriber.AuditTypeURL:  auditLogsQueue.MaxThreads,
	})

	logsPublisher, err := getLogsPublisher(cfg, auditLogsQueue, logger)
	if err != nil {
		return nil, fmt.Errorf("error prepare logs publisher: %s", err.Error())
	}

	metricsService := service.NewMetricsService(storage)

	if cfg.Restore {
		err = restoreMetrics(ctx, cfg.FileStoragePath, storage, metricsService)
		if err != nil {
			return nil, fmt.Errorf("error restore metrics: %s", err.Error())
		}
	}

	logsSaver := service.NewLogsSaverService(fileStorage, cfg.StoreInterval)

	handler := handler.NewMetricsHandler(
		metricsService,
		logsSaver,
		service.NewPingerService(pinger),
		auditLogsQueue,
	)

	return &App{
		Handler:       handler,
		Pool:          pool,
		LogsSaver:     logsSaver,
		LogsPublisher: logsPublisher,
	}, nil
}

func restoreMetrics(
	ctx context.Context,
	fileStoragePath string,
	dbStorage service.MetricsRepository,
	metricsService pkgContracts.MetricsService,
) error {
	fileMetrics, err := file.GetFromFile(fileStoragePath)
	if err != nil {
		return err
	}

	if len(fileMetrics) > 0 {
		existMetrics, err := metricsService.GetGroupedByTypeAndName(ctx)
		if err != nil {
			return err
		}

		var preparedMetrics []model.Metrics
		for _, fileMetric := range fileMetrics {
			if _, ok := existMetrics[fileMetric.MType][fileMetric.Name]; !ok {
				preparedMetrics = append(preparedMetrics, fileMetric)
			}
		}

		dbStorage.AddList(ctx, preparedMetrics)
	}

	return nil
}

func getLogsPublisher(
	cfg *appConfig.Config,
	auditLogs auditlogs.AuditLogs,
	logger *zap.SugaredLogger,
) (publisher.Publisher, error) {
	logsPublisher := publisher.NewLogsPublisherQuery(auditLogs, logger)

	if cfg.AuditFile != "" {
		subscriberAuditFile, err := subscriber.NewAuditFile(cfg.AuditFile)
		if err != nil {
			return nil, err
		}
		logsPublisher.Register(subscriberAuditFile)
	}
	if cfg.AuditURL != "" {
		subscriberAuditURL, err := subscriber.NewAuditURL(cfg.AuditURL)
		if err != nil {
			return nil, err
		}
		logsPublisher.Register(subscriberAuditURL)
	}

	return logsPublisher, nil
}
