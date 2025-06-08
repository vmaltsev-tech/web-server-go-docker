package middleware

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"web-server-go-docker/internal/metrics"

	"github.com/prometheus/client_golang/prometheus/testutil"
)

func TestRequestCounterMiddleware(t *testing.T) {
	counter := 0
	rcm := NewRequestCounterMiddleware(&counter)

	handler := rcm.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/counter", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if counter != 1 {
		t.Errorf("expected counter to be 1, got %d", counter)
	}
}

func TestLoggingMiddlewareRecordsMetrics(t *testing.T) {
	m := metrics.New()
	lm := NewLoggingMiddleware(m)

	handler := lm.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/log", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	counter := testutil.ToFloat64(m.RequestsTotal.WithLabelValues(http.MethodGet, "/log", strconv.Itoa(http.StatusOK)))
	if counter != 1 {
		t.Errorf("expected metric counter 1, got %v", counter)
	}

	if n := testutil.CollectAndCount(m.RequestDuration); n == 0 {
		t.Error("expected request duration metric to be recorded")
	}
}
