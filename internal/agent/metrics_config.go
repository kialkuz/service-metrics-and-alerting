package agent

type MetricsConfig struct {
	ReportInterval int
	PollInterval   int
	RateLimit      int
	CollectAlgo    string
}
