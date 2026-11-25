package audit

import (
	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
)

func GetAuditList(auditConfig appConfig.AuditConfig) (map[string]AuditSaver, error) {
	auditList := make(map[string]AuditSaver)

	if auditConfig.File != "" {
		auditFile, err := NewAuditFile(auditConfig.File)
		if err != nil {
			return nil, err
		}
		auditList[AuditTypeFile] = auditFile
	}

	if auditConfig.URL != "" {
		auditFile, err := NewAuditURL(auditConfig.URL)
		if err != nil {
			return nil, err
		}
		auditList[AuditTypeURL] = auditFile
	}

	return auditList, nil
}
