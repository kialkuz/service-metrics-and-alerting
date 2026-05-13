package agent

type Config struct {
	ReportInterval int
	PollInterval   int
	URL            string
}

func NewConfig() (*Config, error) {
	config, err := GetIncomingParams()
	if err != nil {
		return nil, err
	}

	return &Config{
		ReportInterval: config.ReportInterval,
		PollInterval:   config.PollInterval,
		URL:            "http://" + config.Address,
	}, nil
}
