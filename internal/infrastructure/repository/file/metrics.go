package file

import (
	"bufio"
	"encoding/json"
	"kialkuz/service-metrics-and-alerting/internal/model"
	"os"
)

type Saver struct {
	file   *os.File
	writer *bufio.Writer
}

func (p *Saver) WriteMetric(metric *model.Metrics) error {
	data, err := json.Marshal(&metric)
	if err != nil {
		return err
	}

	data = append(data, '\n')

	_, err = p.file.Write(data)
	return err
}

func (p *Saver) Truncate() error {
	return p.file.Truncate(0)
}

func (p *Saver) Close() error {
	return p.file.Close()
}
