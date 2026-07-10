package agent

type Config struct {
	ReportInterval int
	PollInterval   int
	Key            string
	RateLimit      int
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
		RateLimit:      config.RateLimit,
		URL:            "http://" + config.Address,
	}, nil
}
