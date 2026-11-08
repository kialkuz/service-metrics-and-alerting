package subscriber

import (
	"context"
	"kialkuz/service-metrics-and-alerting/internal/server/event/savelogs"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuditURL_Save_ServerError(t *testing.T) {
	var attempts atomic.Int32

	server := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			attempts.Add(1)
			w.WriteHeader(http.StatusInternalServerError)
		},
	))
	defer server.Close()

	audit, err := NewAuditURL(server.URL)
	require.NoError(t, err)

	err = audit.Save(context.Background(), savelogs.Event{
		TS:        time.Now().Unix(),
		Metrics:   []string{"test"},
		IPAddress: "127.0.0.1",
	})

	require.Error(t, err)
	assert.Equal(t, int32(4), attempts.Load())
}

func TestAuditURL_Save_Success(t *testing.T) {
	expectedEvent := savelogs.Event{}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	audit, err := NewAuditURL(server.URL)
	require.NoError(t, err)

	err = audit.Save(context.Background(), expectedEvent)

	require.NoError(t, err)
}
