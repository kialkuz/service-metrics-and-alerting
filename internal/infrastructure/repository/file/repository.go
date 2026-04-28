package file

import (
	"bufio"
	"encoding/json"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"os"

	_ "github.com/lib/pq"
)

//go:generate go run go.uber.org/mock/mockgen -source=metrics.go -destination=mocks/metrics_mock.go -package=mocks -typed
type MetricsFileRepository interface {
	Truncate() error
	WriteMetric(metric *model.Metrics) error
	Close() error
}

func NewFileStorage(filename string, storeInterval int) (*Saver, error) {
	var file *os.File
	var err error

	if storeInterval != 0 {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	} else {
		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
	}

	if err != nil {
		return nil, err
	}

	return &Saver{file: file, writer: bufio.NewWriter(file)}, nil
}

func GetFromFile(fileStoragePath string) ([]model.Metrics, error) {
	file, _ := os.OpenFile(fileStoragePath, os.O_RDONLY, 0644)

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
