package subscriber

import (
	"context"
	"fmt"
	loggerMocks "kialkuz/service-metrics-and-alerting/internal/logger/mocks"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs/subscriber/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestAuditFile_Save_InvalidFilename(t *testing.T) {
	audit, err := NewAuditFile("")
	require.Nil(t, audit)
	require.Error(t, err)
}

func TestAuditFile_Save_ErrorWriteToLogs(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := loggerMocks.NewMockLogger(ctrl)

	auditSaver := mocks.NewMockAuditSaver(ctrl)
	auditSaver.EXPECT().SetLogger(gomock.Any())
	auditSaver.EXPECT().
		Save(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, event <-chan savelogs.Event) {
			logger.Error(fmt.Errorf("ошибка записи в файл логов аудита"))
		})
	auditSaver.EXPECT().Close().Return(nil)

	var loggerErr error
	logger.EXPECT().
		Error(gomock.Any()).
		Do(func(args ...any) {
			loggerErr = args[0].(error)
		}).
		AnyTimes()
	auditSaver.SetLogger(logger)

	events := make(chan savelogs.Event, 1)

	events <- savelogs.Event{
		TS:        time.Now().Unix(),
		Metrics:   []string{"test"},
		IPAddress: "127.0.0.1",
	}

	close(events)

	auditSaver.Save(context.TODO(), events)

	require.Error(t, loggerErr)

	err := auditSaver.Close()
	require.NoError(t, err)
}

func TestAuditFile_Save_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := loggerMocks.NewMockLogger(ctrl)

	auditSaver := mocks.NewMockAuditSaver(ctrl)
	auditSaver.EXPECT().SetLogger(gomock.Any())
	auditSaver.EXPECT().Save(gomock.Any(), gomock.Any())
	auditSaver.EXPECT().Close().Return(nil)

	var loggerErr error
	logger.EXPECT().
		Error(gomock.Any()).
		Do(func(args ...any) {
			loggerErr = args[0].(error)
		}).
		AnyTimes()
	auditSaver.SetLogger(logger)

	events := make(chan savelogs.Event, 1)

	events <- savelogs.Event{
		TS:        time.Now().Unix(),
		Metrics:   []string{"test"},
		IPAddress: "127.0.0.1",
	}

	close(events)

	auditSaver.Save(context.TODO(), events)

	require.NoError(t, loggerErr)

	err := auditSaver.Close()
	require.NoError(t, err)
}
