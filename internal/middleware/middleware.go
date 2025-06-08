package middleware

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"web-server-go-docker/internal/metrics"
)

// statusResponseWriter оборачивает ResponseWriter для захвата HTTP статуса
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

// Middleware представляет интерфейс для middleware
type Middleware interface {
	Handler(next http.Handler) http.Handler
}

// LoggingMiddleware логирует HTTP запросы и собирает метрики
type LoggingMiddleware struct {
	metrics *metrics.Metrics
}

// NewLoggingMiddleware создает новый LoggingMiddleware
func NewLoggingMiddleware(m *metrics.Metrics) *LoggingMiddleware {
	return &LoggingMiddleware{
		metrics: m,
	}
}

// Handler возвращает middleware handler для логирования
func (lm *LoggingMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Оборачиваем ResponseWriter для захвата статуса
		wrapped := &statusResponseWriter{ResponseWriter: w, statusCode: 200}

		next.ServeHTTP(wrapped, r)

		duration := time.Since(start)

		// Логируем
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, wrapped.statusCode, duration)

		// Собираем метрики если они доступны
		if lm.metrics != nil {
			lm.metrics.RecordRequest(
				r.Method,
				r.URL.Path,
				strconv.Itoa(wrapped.statusCode),
				duration,
			)
		}
	})
}

// SecurityMiddleware добавляет security headers
type SecurityMiddleware struct{}

// NewSecurityMiddleware создает новый SecurityMiddleware
func NewSecurityMiddleware() *SecurityMiddleware {
	return &SecurityMiddleware{}
}

// Handler возвращает middleware handler для security headers
func (sm *SecurityMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		next.ServeHTTP(w, r)
	})
}

// RequestCounterMiddleware подсчитывает количество запросов
type RequestCounterMiddleware struct {
	counter *int
}

// NewRequestCounterMiddleware создает новый RequestCounterMiddleware
func NewRequestCounterMiddleware(counter *int) *RequestCounterMiddleware {
	return &RequestCounterMiddleware{
		counter: counter,
	}
}

// Handler возвращает middleware handler для подсчета запросов
func (rcm *RequestCounterMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if rcm.counter != nil {
			(*rcm.counter)++
		}
		next.ServeHTTP(w, r)
	})
}

// Chain объединяет несколько middleware в цепочку
func Chain(middlewares ...Middleware) func(http.Handler) http.Handler {
	return func(final http.Handler) http.Handler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			final = middlewares[i].Handler(final)
		}
		return final
	}
}
