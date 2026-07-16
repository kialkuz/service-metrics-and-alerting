package file

import (
	"bufio"
	"encoding/json"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"os"
)

func GetFromFile(fileStoragePath string) ([]model.Metrics, error) {
	file, err := os.OpenFile(fileStoragePath+"/metrics", os.O_CREATE|os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

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
	}

	return metricsList, nil
}
