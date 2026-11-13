package subscriber

import (
	"context"
	"errors"
	loggerMocks "kialkuz/service-metrics-and-alerting/internal/logger/mocks"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	loggerRetrier "kialkuz/service-metrics-and-alerting/internal/service/retrier/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuditFile_Save_InvalidUrl(t *testing.T) {
	audit, err := NewAuditURL("test.ru")
	require.Nil(t, audit)
	require.Error(t, err)
}

func TestAuditURL_Save_Error(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := loggerMocks.NewMockLogger(ctrl)
	retrierMock := loggerRetrier.NewMockRetrier(ctrl)

	var loggerErr error

	logger.EXPECT().
		Error(gomock.Any()).
		Do(func(args ...any) {
			loggerErr = args[0].(error)
		}).
		AnyTimes()

	retrierMock.EXPECT().
		WithRetry(gomock.Any(), gomock.Any()).
		Return(errors.New("test error"))

	auditURL := &AuditURL{
		Audit:    Audit{AuditType: AuditTypeURL},
		savePath: "http://localhost:8080/audit",
		retrier:  retrierMock,
	}

	auditURL.SetLogger(logger)

	events := make(chan savelogs.Event, 1)

	events <- savelogs.Event{
		TS:        time.Now().Unix(),
		Metrics:   []string{"test"},
		IPAddress: "127.0.0.1",
	}

	close(events)

	auditURL.Save(context.TODO(), events)

	require.Error(t, loggerErr)
	require.EqualError(t, loggerErr, "test error")
}

func TestAuditUrl_Save_Success(t *testing.T) {
	ctrl := gomock.NewController(t)

	logger := loggerMocks.NewMockLogger(ctrl)
	retrierMock := loggerRetrier.NewMockRetrier(ctrl)

	var loggerErr error

	logger.EXPECT().
		Error(gomock.Any()).
		Do(func(args ...any) {
			loggerErr = args[0].(error)
		}).
		AnyTimes()

	retrierMock.EXPECT().
		WithRetry(gomock.Any(), gomock.Any()).
		Return(nil)

	auditURL := &AuditURL{
		Audit:    Audit{AuditType: AuditTypeURL},
		savePath: "http://localhost:8080/audit",
		retrier:  retrierMock,
	}

	auditURL.SetLogger(logger)

	events := make(chan savelogs.Event, 1)

	events <- savelogs.Event{
		TS:        time.Now().Unix(),
		Metrics:   []string{"test"},
		IPAddress: "127.0.0.1",
	}

	close(events)

	auditURL.Save(context.TODO(), events)

	require.NoError(t, loggerErr)
}
