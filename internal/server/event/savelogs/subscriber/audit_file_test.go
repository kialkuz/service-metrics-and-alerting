package subscriber

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAuditFile_Save_InvalidFilename(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	audit, err := NewAuditFile("")
	require.Nil(t, audit)
	require.Error(t, err)
}

func TestAuditFile_Save_Success(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	audit, err := NewAuditFile("./test.log")
	require.NoError(t, err)

	err = audit.Save(context.Background(), savelogs.Event{
		TS:        time.Now().Unix(),
		Metrics:   []string{"test"},
		IPAddress: "127.0.0.1",
	})
	require.NoError(t, err)

	err = audit.Close()
	require.NoError(t, err)
}
