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

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	Handler       *handler.MetricsHandler
	Pool          *pgxpool.Pool
	FileService   *service.FileService
	LogsPublisher publisher.Publisher
	AuditLogsList *service.AuditLogsList
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

	auditLogsList := service.NewAuditLogsList()

	logsPublisher, err := getLogsPublisher(cfg, auditLogsList, logger)
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

	fileService := service.NewFileService(fileStorage, cfg.StoreInterval)

	handler := handler.NewMetricsHandler(
		metricsService,
		fileService,
		service.NewPingerService(pinger),
		logsPublisher,
	)

	return &App{
		Handler:       handler,
		Pool:          pool,
		FileService:   fileService,
		AuditLogsList: auditLogsList,
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
	auditLogsList *service.AuditLogsList,
	logger *zap.SugaredLogger,
) (publisher.Publisher, error) {
	logsPublisher := publisher.NewLogsPublisher(auditLogsList, logger)

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
