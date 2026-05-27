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
	"kialkuz/service-metrics-and-alerting/internal/server/handler"

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

	if cfg.Restore {
		err = restoreMetrics(ctx, cfg.FileStoragePath, storage)
		if err != nil {
			return nil, fmt.Errorf("error restore metrics: %s", err.Error())
		}
	}

	fileService := service.NewFileService(fileStorage, cfg.StoreInterval)

	handler := handler.NewMetricsHandler(
		service.NewMetricsService(storage),
		fileService,
		service.NewPingerService(pinger),
	)

	return &App{
		Handler:     handler,
		Pool:        pool,
		FileService: fileService,
	}, nil
}

func restoreMetrics(ctx context.Context, fileStoragePath string, dbStorage service.MetricsRepository) error {
	metrics, err := file.GetFromFile(fileStoragePath)
	if err != nil {
		return err
	}

	if len(metrics) > 0 {
		dbStorage.AddList(ctx, metrics)
	}

	return nil
}
