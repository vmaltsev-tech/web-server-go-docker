package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"web-server-go-docker/internal/config"
	"web-server-go-docker/internal/handlers"
	"web-server-go-docker/internal/metrics"
	"web-server-go-docker/internal/middleware"
)

// Server представляет HTTP сервер с зависимостями
type Server struct {
	config       *config.Config
	metrics      *metrics.Metrics
	handler      *handlers.Handler
	requestCount int
	httpServer   *http.Server
}

// New создает новый сервер с зависимостями
func New(cfg *config.Config) (*Server, error) {
	var m *metrics.Metrics
	if cfg.Metrics.Enabled {
		m = metrics.New()
	}

	s := &Server{
		config:       cfg,
		metrics:      m,
		requestCount: 0,
	}

	s.handler = handlers.New(cfg, m, &s.requestCount)
	s.setupRoutes()

	return s, nil
}

// setupRoutes настраивает маршруты и middleware
func (s *Server) setupRoutes() {
	mux := http.NewServeMux()
	
	// Регистрируем маршруты
	mux.HandleFunc("/", s.handler.Info)
	mux.HandleFunc("/health", s.handler.Health)
	mux.HandleFunc("/metrics", s.handler.Metrics)
	
	if s.config.Metrics.Enabled && s.metrics != nil {
		mux.HandleFunc(s.config.Metrics.Path, s.handler.PrometheusMetrics)
	}

	// Настраиваем middleware
	var middlewares []middleware.Middleware
	
	middlewares = append(middlewares, middleware.NewSecurityMiddleware())
	middlewares = append(middlewares, middleware.NewRequestCounterMiddleware(&s.requestCount))
	
	if s.metrics != nil {
		middlewares = append(middlewares, middleware.NewLoggingMiddleware(s.metrics))
	}

	// Применяем middleware chain
	handler := middleware.Chain(middlewares...)(mux)

	s.httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%s", s.config.Server.Port),
		Handler:      handler,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}
}

// Start запускает сервер
func (s *Server) Start() error {
	// Канал для graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Запуск сервера в отдельной горутине
	go func() {
		log.Printf("Starting server on port %s", s.config.Server.Port)
		log.Printf("Environment: %s", s.config.App.Environment)
		log.Printf("Available endpoints:")
		log.Printf("  GET / - Server info")
		log.Printf("  GET /health - Health check")
		log.Printf("  GET /metrics - Server metrics (JSON)")
		
		if s.config.Metrics.Enabled {
			log.Printf("  GET %s - Prometheus metrics", s.config.Metrics.Path)
		}

		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on port %s: %v", s.config.Server.Port, err)
		}
	}()

	// Ожидание сигнала завершения
	<-quit
	log.Println("Shutting down server...")

	return s.Shutdown()
}

// Shutdown выполняет graceful shutdown сервера
func (s *Server) Shutdown() error {
	// Graceful shutdown с таймаутом 10 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}

	log.Println("Server exited gracefully")
	return nil
}

// GetRequestCount возвращает количество обработанных запросов
func (s *Server) GetRequestCount() int {
	return s.requestCount
} 