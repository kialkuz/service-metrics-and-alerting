package audit

import (
	"fmt"
	"reflect"

	appConfig "kialkuz/service-metrics-and-alerting/internal/config/server"
)

func GetAuditList(auditConfig appConfig.AuditConfig) (map[string]AuditSaver, error) {
	auditList := make(map[string]AuditSaver)

	config := reflect.ValueOf(auditConfig)
	configType := config.Type()

	for i := 0; i < config.NumField(); i++ {
		value := config.Field(i)

		if value.IsZero() {
			continue
		}

		name := configType.Field(i).Name

		constructor, ok := GetAuditSubscriberByName(name)
		if !ok {
			return nil, fmt.Errorf("constructor for audit %s not found", name)
		}

		audit, err := constructor(value.String())
		if err != nil {
			return nil, err
		}

		auditList[name] = audit
	}

	return auditList, nil
}
