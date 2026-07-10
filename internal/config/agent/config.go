package agent

type Config struct {
	ReportInterval int
	PollInterval   int
	Key            string
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
		Key:            config.Key,
		URL:            "http://" + config.Address,
	}, nil
}
