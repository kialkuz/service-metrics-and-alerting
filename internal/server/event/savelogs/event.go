package savelogs

const (
	AuditFile = "file"
	AuditUrl  = "url"
)

type Event struct {
	Ts        int64    `json:"ts"`
	Metrics   []string `json:"metrics"`
	IpAddress string   `json:"ipAddress"`
}
