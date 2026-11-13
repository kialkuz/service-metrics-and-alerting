package logger

//go:generate go run go.uber.org/mock/mockgen -source=logger.go -destination=mocks/logger_mock.go -package=mocks -typed
type Logger interface {
	Error(args ...any)
}
