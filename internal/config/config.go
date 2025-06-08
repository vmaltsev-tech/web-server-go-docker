package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config представляет конфигурацию приложения
type Config struct {
	Server   ServerConfig
	App      AppConfig
	Metrics  MetricsConfig
	Logging  LoggingConfig
}

// ServerConfig содержит настройки HTTP сервера
type ServerConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// AppConfig содержит настройки приложения
type AppConfig struct {
	Environment string
	Version     string
}

// MetricsConfig содержит настройки метрик
type MetricsConfig struct {
	Enabled bool
	Path    string
}

// LoggingConfig содержит настройки логирования
type LoggingConfig struct {
	Level  string
	Format string
}

// Load загружает конфигурацию из переменных окружения с валидацией
func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port:         getEnv("PORT", "8080"),
			ReadTimeout:  getDurationEnv("READ_TIMEOUT", 15*time.Second),
			WriteTimeout: getDurationEnv("WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:  getDurationEnv("IDLE_TIMEOUT", 60*time.Second),
		},
		App: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			Version:     getEnv("APP_VERSION", "1.0.0"),
		},
		Metrics: MetricsConfig{
			Enabled: getBoolEnv("METRICS_ENABLED", true),
			Path:    getEnv("METRICS_PATH", "/prometheus"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}

	if err := config.validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

// validate проверяет корректность конфигурации
func (c *Config) validate() error {
	// Валидация порта
	if port, err := strconv.Atoi(c.Server.Port); err != nil || port < 1 || port > 65535 {
		return fmt.Errorf("invalid port: %s", c.Server.Port)
	}

	// Валидация окружения
	validEnvs := map[string]bool{
		"development": true,
		"staging":     true,
		"production":  true,
		"test":        true,
	}
	if !validEnvs[c.App.Environment] {
		return fmt.Errorf("invalid environment: %s", c.App.Environment)
	}

	// Валидация log level
	validLogLevels := map[string]bool{
		"debug": true,
		"info":  true,
		"warn":  true,
		"error": true,
		"fatal": true,
	}
	if !validLogLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	// Валидация log format
	validLogFormats := map[string]bool{
		"json": true,
		"text": true,
	}
	if !validLogFormats[c.Logging.Format] {
		return fmt.Errorf("invalid log format: %s", c.Logging.Format)
	}

	return nil
}

// IsProduction возвращает true если окружение production
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// IsDevelopment возвращает true если окружение development
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getDurationEnv возвращает duration из переменной окружения или значение по умолчанию
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

// getBoolEnv возвращает bool из переменной окружения или значение по умолчанию
func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
} 