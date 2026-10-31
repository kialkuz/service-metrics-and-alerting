package savelogs

const (
	AuditFile = "file"
	AuditURL  = "url"
)

type Event struct {
	TS        int64    `json:"ts"`
	Metrics   []string `json:"metrics"`
	IPAddress string   `json:"ip_address"`
}
