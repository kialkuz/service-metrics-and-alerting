package agent

type Config struct {
	ReportInterval int
	PollInterval   int
	URL            string
}

func NewConfig() (*Config, error) {
	incomingParams, err := GetIncomingParams()
	if err != nil {
		return nil, err
	}

	return &Config{
		ReportInterval: incomingParams.ReportInterval,
		PollInterval:   incomingParams.PollInterval,
		URL:            "http://" + incomingParams.Address,
	}, nil
}
