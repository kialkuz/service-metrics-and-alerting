package file

import (
	"bufio"
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"os"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsFileRepository interface {
	CreateTemp() error
	WriteMetric(metric *model.Metrics) error
	SaveOriginal() error
}

func NewFileStorage(config *server.Config) (*Saver, error) {
	if err := os.MkdirAll(config.FileStoragePath, 0644); err != nil {
		return nil, fmt.Errorf("cannot create folder: %w", err)
	}

	return &Saver{config: config}, nil
}

func GetFromFile(fileStoragePath string) ([]model.Metrics, error) {
	file, err := os.OpenFile(fileStoragePath+"/metrics", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	var metricsList []model.Metrics

	if file != nil {
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			var metrics model.Metrics

			data := scanner.Bytes()

			err := json.Unmarshal(data, &metrics)
			if err != nil {
				return nil, err
			}

			metricsList = append(metricsList, metrics)
		}

		file.Close()
	}

	return metricsList, nil
}
