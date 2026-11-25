package errorclassifier

// ErrorClassification тип для классификации ошибок
type ErrorClassification int

const (
	NonRetriable ErrorClassification = iota
	Retriable
)

type ErrorClassifier interface {
	Classify(err error) ErrorClassification
}
