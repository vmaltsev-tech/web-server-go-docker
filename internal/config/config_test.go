
package config

import (
	"os"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	// Сохраняем оригинальные значения
	originalEnvs := map[string]string{
		"PORT":        os.Getenv("PORT"),
		"ENVIRONMENT": os.Getenv("ENVIRONMENT"),
		"APP_VERSION": os.Getenv("APP_VERSION"),
		"LOG_LEVEL":   os.Getenv("LOG_LEVEL"),
	}

	// Очищаем переменные окружения после теста
	defer func() {
		for key, value := range originalEnvs {
			if value == "" {
				os.Unsetenv(key)
			} else {
				os.Setenv(key, value)
			}
		}
	}()

	tests := []struct {
		name    string
		envVars map[string]string
		wantErr bool
	}{
		{
			name: "default values",
			envVars: map[string]string{
				"PORT":        "",
				"ENVIRONMENT": "",
				"APP_VERSION": "",
				"LOG_LEVEL":   "",
			},
			wantErr: false,
		},
		{
			name: "valid custom values",
			envVars: map[string]string{
				"PORT":        "3000",
				"ENVIRONMENT": "production",
				"APP_VERSION": "2.0.0",
				"LOG_LEVEL":   "debug",
			},
			wantErr: false,
		},
		{
			name: "invalid port",
			envVars: map[string]string{
				"PORT": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid environment",
			envVars: map[string]string{
				"ENVIRONMENT": "invalid",
			},
			wantErr: true,
		},
		{
			name: "invalid log level",
			envVars: map[string]string{
				"LOG_LEVEL": "invalid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Устанавливаем переменные окружения для теста
			for key, value := range tt.envVars {
				if value == "" {
					os.Unsetenv(key)
				} else {
					os.Setenv(key, value)
				}
			}

			cfg, err := Load()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Load() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Load() unexpected error: %v", err)
				return
			}

			// Проверяем значения по умолчанию
			if tt.name == "default values" {
				if cfg.Server.Port != "8080" {
					t.Errorf("Expected default port 8080, got %s", cfg.Server.Port)
				}
				if cfg.App.Environment != "development" {
					t.Errorf("Expected default environment development, got %s", cfg.App.Environment)
				}
				if cfg.App.Version != "1.0.0" {
					t.Errorf("Expected default version 1.0.0, got %s", cfg.App.Version)
				}
			}
		})
	}
}

func TestConfigMethods(t *testing.T) {
	cfg := &Config{
		App: AppConfig{
			Environment: "production",
		},
	}

	if !cfg.IsProduction() {
		t.Error("Expected IsProduction() to return true for production environment")
	}

	if cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return false for production environment")
	}

	cfg.App.Environment = "development"

	if cfg.IsProduction() {
		t.Error("Expected IsProduction() to return false for development environment")
	}

	if !cfg.IsDevelopment() {
		t.Error("Expected IsDevelopment() to return true for development environment")
	}
}

func TestGetDurationEnv(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{
			name:         "valid duration",
			envValue:     "30s",
			defaultValue: 15 * time.Second,
			expected:     30 * time.Second,
		},
		{
			name:         "invalid duration",
			envValue:     "invalid",
			defaultValue: 15 * time.Second,
			expected:     15 * time.Second,
		},
		{
			name:         "empty value",
			envValue:     "",
			defaultValue: 15 * time.Second,
			expected:     15 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_DURATION"
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
			} else {
				os.Unsetenv(key)
			}
			defer os.Unsetenv(key)

			result := getDurationEnv(key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestGetBoolEnv(t *testing.T) {
	tests := []struct {
		name         string
		envValue     string
		defaultValue bool
		expected     bool
	}{
		{
			name:         "true value",
			envValue:     "true",
			defaultValue: false,
			expected:     true,
		},
		{
			name:         "false value",
			envValue:     "false",
			defaultValue: true,
			expected:     false,
		},
		{
			name:         "invalid value",
			envValue:     "invalid",
			defaultValue: true,
			expected:     true,
		},
		{
			name:         "empty value",
			envValue:     "",
			defaultValue: true,
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			key := "TEST_BOOL"
			if tt.envValue != "" {
				os.Setenv(key, tt.envValue)
			} else {
				os.Unsetenv(key)
			}
			defer os.Unsetenv(key)

			result := getBoolEnv(key, tt.defaultValue)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
