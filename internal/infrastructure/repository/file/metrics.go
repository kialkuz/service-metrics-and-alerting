package file

import (
	"encoding/json"
	"fmt"
	"kialkuz/service-metrics-and-alerting/internal/config/server"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"os"
)

type Saver struct {
	file   *os.File
	config *server.Config
}

func (p *Saver) CreateTemp() error {
	file, err := os.CreateTemp(p.config.FileStoragePath, "*")
	if err != nil {
		return err
	}

	p.file = file

	return nil
}

func (p *Saver) WriteMetric(metric *model.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = p.file.Write(data)
	if err != nil {
		return err
	}
	return err
}

func (p *Saver) SaveOriginal() error {
	if p.file != nil {
		p.file.Close()

		err := os.Rename(p.file.Name(), p.config.FileStoragePath+"/metrics")
		if err != nil {
			return fmt.Errorf("error rename to original file: %s", err.Error())
		}
	}

	return nil
}
