package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics содержит все Prometheus метрики
type Metrics struct {
	RequestsTotal    *prometheus.CounterVec
	RequestDuration  *prometheus.HistogramVec
	ServerUptime     *prometheus.GaugeVec
	startTime        time.Time
	registry         *prometheus.Registry
}

// New создает новый экземпляр метрик
func New() *Metrics {
	// Создаем новый registry для избежания конфликтов в тестах
	registry := prometheus.NewRegistry()
	
	requestsTotal := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests.",
		},
		[]string{"method", "endpoint", "status"},
	)
	
	requestDuration := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)
	
	serverUptime := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "server_uptime_seconds",
			Help: "Server uptime in seconds.",
		},
		nil,
	)

	m := &Metrics{
		RequestsTotal:   requestsTotal,
		RequestDuration: requestDuration,
		ServerUptime:    serverUptime,
		startTime:       time.Now(),
		registry:        registry,
	}

	// Регистрируем метрики в нашем registry
	registry.MustRegister(requestsTotal)
	registry.MustRegister(requestDuration)
	registry.MustRegister(serverUptime)

	return m
}

// RecordRequest записывает метрики HTTP запроса
func (m *Metrics) RecordRequest(method, endpoint, status string, duration time.Duration) {
	m.RequestsTotal.WithLabelValues(method, endpoint, status).Inc()
	m.RequestDuration.WithLabelValues(method, endpoint).Observe(duration.Seconds())
}

// UpdateUptime обновляет метрику uptime
func (m *Metrics) UpdateUptime() {
	m.ServerUptime.WithLabelValues().Set(time.Since(m.startTime).Seconds())
}

// GetStartTime возвращает время запуска сервера
func (m *Metrics) GetStartTime() time.Time {
	return m.startTime
}

// Handler возвращает HTTP handler для Prometheus метрик
func (m *Metrics) Handler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})
} 