package collectalgo

type CollectConfig struct {
	CountAdditionalReplaySendMetrics int
	AdditionalReplaySendInterval     map[int]int
	RateLimit                        int
}
