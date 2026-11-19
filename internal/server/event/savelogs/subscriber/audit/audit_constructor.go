package audit

type AuditConstructor func(string) (AuditSaver, error)

var auditRegistry = make(map[string]AuditConstructor)

func RegisterAuditSubscriber(name string, constructor AuditConstructor) {
	auditRegistry[name] = constructor
}

func GetAuditRegistry() map[string]AuditConstructor {
	return auditRegistry
}

func GetAuditSubscriberByName(name string) (AuditConstructor, bool) {
	subscriber, ok := auditRegistry[name]

	return subscriber, ok
}
