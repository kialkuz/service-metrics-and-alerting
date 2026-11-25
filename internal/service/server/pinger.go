package server

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/infrastructure/repository/db"
	pkgErrors "kialkuz/service-metrics-and-alerting/pkg/errors"
)

type PingerService struct {
	repository *db.MemStorage
}

func NewPingerService(repository *db.MemStorage) *PingerService {
	return &PingerService{repository: repository}
}

func (s *PingerService) Ping(ctx context.Context) error {
	if s.repository == nil {
		return pkgErrors.ErrNotInitDB
	}

	err := s.repository.Ping(ctx)
	if err != nil {
		return err
	}

	return nil
}
