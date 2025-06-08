package integration

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"web-server-go-docker/internal/config"
	"web-server-go-docker/internal/models"
	"web-server-go-docker/internal/server"
)

func TestServerIntegration(t *testing.T) {
	// Создаем тестовую конфигурацию
	cfg := &config.Config{
		Server: config.ServerConfig{
			Port:         "0", // Используем случайный порт
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 5 * time.Second,
			IdleTimeout:  10 * time.Second,
		},
		App: config.AppConfig{
			Environment: "test",
			Version:     "1.0.0-test",
		},
		Metrics: config.MetricsConfig{
			Enabled: true,
			Path:    "/prometheus",
		},
		Logging: config.LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}

	// Создаем сервер
	srv, err := server.New(cfg)
	if err != nil {
		t.Fatalf("Failed to create server: %v", err)
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		if err := srv.Start(); err != nil {
			t.Logf("Server error: %v", err)
		}
	}()

	// Даем серверу время на запуск
	time.Sleep(100 * time.Millisecond)

	baseURL := fmt.Sprintf("http://localhost:%s", cfg.Server.Port)

	// Тестируем endpoints
	t.Run("health endpoint", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var healthResp models.HealthResponse
		if err := json.NewDecoder(resp.Body).Decode(&healthResp); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if healthResp.Status != "OK" {
			t.Errorf("Expected status 'OK', got '%s'", healthResp.Status)
		}

		if healthResp.Version != "1.0.0-test" {
			t.Errorf("Expected version '1.0.0-test', got '%s'", healthResp.Version)
		}
	})

	t.Run("info endpoint", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var infoResp models.InfoResponse
		if err := json.NewDecoder(resp.Body).Decode(&infoResp); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		expectedMessage := "DevOps Portfolio 2025 - Go Web Server"
		if infoResp.Message != expectedMessage {
			t.Errorf("Expected message '%s', got '%s'", expectedMessage, infoResp.Message)
		}

		if infoResp.Environment != "test" {
			t.Errorf("Expected environment 'test', got '%s'", infoResp.Environment)
		}
	})

	t.Run("metrics endpoint", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/metrics")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		var metricsResp models.MetricsResponse
		if err := json.NewDecoder(resp.Body).Decode(&metricsResp); err != nil {
			t.Errorf("Failed to decode response: %v", err)
		}

		if metricsResp.RequestCount < 0 {
			t.Errorf("Expected non-negative request count, got %d", metricsResp.RequestCount)
		}

		if metricsResp.Uptime == "" {
			t.Error("Expected uptime to be set")
		}

		if metricsResp.StartTime == "" {
			t.Error("Expected start time to be set")
		}
	})

	t.Run("prometheus endpoint", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/prometheus")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Errorf("Expected status 200, got %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType != "text/plain; version=0.0.4; charset=utf-8" {
			t.Errorf("Expected Prometheus content type, got '%s'", contentType)
		}
	})

	t.Run("security headers", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			t.Fatalf("Failed to make request: %v", err)
		}
		defer resp.Body.Close()

		expectedHeaders := map[string]string{
			"X-Content-Type-Options": "nosniff",
			"X-Frame-Options":        "DENY",
			"X-XSS-Protection":       "1; mode=block",
			"Referrer-Policy":        "strict-origin-when-cross-origin",
		}

		for header, expectedValue := range expectedHeaders {
			actualValue := resp.Header.Get(header)
			if actualValue != expectedValue {
				t.Errorf("Expected header %s: %s, got %s", header, expectedValue, actualValue)
			}
		}
	})

	// Останавливаем сервер
	if err := srv.Shutdown(); err != nil {
		t.Errorf("Failed to shutdown server: %v", err)
	}
} 