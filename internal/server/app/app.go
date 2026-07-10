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

	"github.com/jackc/pgx/v5/pgxpool"

	service "kialkuz/service-metrics-and-alerting/internal/service/server"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type App struct {
	Handler     *handler.MetricsHandler
	Pool        *pgxpool.Pool
	FileService *service.FileService
}

func NewApp(ctx context.Context, cfg *appConfig.Config) (*App, error) {
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
	)

	return &App{
		Handler:     handler,
		Pool:        pool,
		FileService: fileService,
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
