
package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"web-server-go-docker/internal/config"
	"web-server-go-docker/internal/metrics"
	"web-server-go-docker/internal/models"
)

// Handler содержит зависимости для обработчиков
type Handler struct {
	config        *config.Config
	metrics       *metrics.Metrics
	requestCount  *int
}

// New создает новый Handler с зависимостями
func New(cfg *config.Config, m *metrics.Metrics, requestCount *int) *Handler {
	return &Handler{
		config:       cfg,
		metrics:      m,
		requestCount: requestCount,
	}
}

// Health обрабатывает health check запросы
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.HealthResponse{
		Status:    "OK",
		Timestamp: time.Now().Format(time.RFC3339),
		Version:   h.config.App.Version,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding health response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Info обрабатывает info запросы
func (h *Handler) Info(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	response := models.InfoResponse{
		Message:     "DevOps Portfolio 2025 - Go Web Server",
		Environment: h.config.App.Environment,
		Port:        h.config.Server.Port,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding info response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// Metrics обрабатывает metrics запросы (JSON формат)
func (h *Handler) Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	
	var requestCount int
	if h.requestCount != nil {
		requestCount = *h.requestCount
	}

	var startTime time.Time
	var uptime string
	if h.metrics != nil {
		startTime = h.metrics.GetStartTime()
		uptime = time.Since(startTime).String()
	} else {
		startTime = time.Now()
		uptime = "0s"
	}

	response := models.MetricsResponse{
		RequestCount: requestCount,
		Uptime:       uptime,
		StartTime:    startTime.Format(time.RFC3339),
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding metrics response: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}

// PrometheusMetrics обрабатывает Prometheus metrics запросы
func (h *Handler) PrometheusMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.metrics == nil {
		http.Error(w, "Metrics not available", http.StatusServiceUnavailable)
		return
	}

	// Обновляем uptime перед отдачей метрик
	h.metrics.UpdateUptime()

	// Используем стандартный Prometheus handler
	h.metrics.Handler().ServeHTTP(w, r)
}
