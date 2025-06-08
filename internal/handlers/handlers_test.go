package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"web-server-go-docker/internal/config"
	"web-server-go-docker/internal/metrics"
	"web-server-go-docker/internal/models"
)

func TestHandler_Health(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Version: "1.0.0",
		},
	}
	
	requestCount := 0
	h := New(cfg, nil, &requestCount)

	tests := []struct {
		name           string
		method         string
		expectedStatus int
		expectedBody   bool
	}{
		{
			name:           "GET request should return 200",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:           "POST request should return 405",
			method:         http.MethodPost,
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/health", nil)
			w := httptest.NewRecorder()

			h.Health(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody {
				var response models.HealthResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				if response.Status != "OK" {
					t.Errorf("Expected status 'OK', got '%s'", response.Status)
				}

				if response.Version != "1.0.0" {
					t.Errorf("Expected version '1.0.0', got '%s'", response.Version)
				}

				contentType := w.Header().Get("Content-Type")
				if contentType != "application/json" {
					t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
				}
			}
		})
	}
}

func TestHandler_Info(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Environment: "test",
		},
		Server: config.ServerConfig{
			Port: "8080",
		},
	}
	
	requestCount := 0
	h := New(cfg, nil, &requestCount)

	tests := []struct {
		name           string
		method         string
		path           string
		expectedStatus int
		expectedBody   bool
	}{
		{
			name:           "GET request to root should return 200",
			method:         http.MethodGet,
			path:           "/",
			expectedStatus: http.StatusOK,
			expectedBody:   true,
		},
		{
			name:           "GET request to other path should return 404",
			method:         http.MethodGet,
			path:           "/other",
			expectedStatus: http.StatusNotFound,
			expectedBody:   false,
		},
		{
			name:           "POST request should return 405",
			method:         http.MethodPost,
			path:           "/",
			expectedStatus: http.StatusMethodNotAllowed,
			expectedBody:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			w := httptest.NewRecorder()

			h.Info(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if tt.expectedBody {
				var response models.InfoResponse
				if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
					t.Errorf("Failed to unmarshal response: %v", err)
				}

				expectedMessage := "DevOps Portfolio 2025 - Go Web Server"
				if response.Message != expectedMessage {
					t.Errorf("Expected message '%s', got '%s'", expectedMessage, response.Message)
				}

				if response.Environment != "test" {
					t.Errorf("Expected environment 'test', got '%s'", response.Environment)
				}

				if response.Port != "8080" {
					t.Errorf("Expected port '8080', got '%s'", response.Port)
				}
			}
		})
	}
}

func TestHandler_Metrics(t *testing.T) {
	cfg := &config.Config{}
	requestCount := 5
	m := metrics.New()
	h := New(cfg, m, &requestCount)

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	h.Metrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.MetricsResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.RequestCount != 5 {
		t.Errorf("Expected request count 5, got %d", response.RequestCount)
	}

	if response.Uptime == "" {
		t.Error("Expected uptime to be set")
	}

	if response.StartTime == "" {
		t.Error("Expected start time to be set")
	}
}

func TestHandler_PrometheusMetrics(t *testing.T) {
	cfg := &config.Config{}
	requestCount := 0

	tests := []struct {
		name           string
		metrics        *metrics.Metrics
		expectedStatus int
	}{
		{
			name:           "with metrics enabled",
			metrics:        metrics.New(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "with metrics disabled",
			metrics:        nil,
			expectedStatus: http.StatusServiceUnavailable,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := New(cfg, tt.metrics, &requestCount)

			req := httptest.NewRequest(http.MethodGet, "/prometheus", nil)
			w := httptest.NewRecorder()

			h.PrometheusMetrics(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	cfg := &config.Config{
		App: config.AppConfig{
			Version: "1.0.0",
		},
	}
	requestCount := 0
	h := New(cfg, nil, &requestCount)

	endpoints := []struct {
		name    string
		handler func(http.ResponseWriter, *http.Request)
		path    string
	}{
		{"health", h.Health, "/health"},
		{"metrics", h.Metrics, "/metrics"},
		{"prometheus", h.PrometheusMetrics, "/prometheus"},
	}

	methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

	for _, endpoint := range endpoints {
		for _, method := range methods {
			t.Run(endpoint.name+"_"+method, func(t *testing.T) {
				req := httptest.NewRequest(method, endpoint.path, nil)
				w := httptest.NewRecorder()

				endpoint.handler(w, req)

				if w.Code != http.StatusMethodNotAllowed {
					t.Errorf("Expected status %d for %s %s, got %d", 
						http.StatusMethodNotAllowed, method, endpoint.path, w.Code)
				}
			})
		}
	}
} 