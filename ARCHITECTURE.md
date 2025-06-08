# Архитектура проекта

## Обзор

Проект был полностью рефакторен с использованием принципов чистой архитектуры и лучших практик Go. Новая структура обеспечивает лучшую модульность, тестируемость и поддерживаемость кода.

## Структура проекта

```
web-server-go-docker/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── config/
│   │   ├── config.go            # Конфигурация приложения
│   │   └── config_test.go       # Unit тесты конфигурации
│   ├── handlers/
│   │   ├── handlers.go          # HTTP обработчики
│   │   └── handlers_test.go     # Unit тесты обработчиков
│   ├── metrics/
│   │   └── prometheus.go        # Prometheus метрики
│   ├── middleware/
│   │   └── middleware.go        # HTTP middleware
│   ├── models/
│   │   └── responses.go         # Модели ответов
│   └── server/
│       └── server.go            # HTTP сервер
├── test/
│   └── integration/
│       └── server_test.go       # Integration тесты
├── monitoring/                  # Конфигурация мониторинга
├── go.mod
├── go.sum
├── Dockerfile
├── docker-compose.yaml
└── Makefile
```

## Принципы архитектуры

### 1. Разделение ответственности (Separation of Concerns)

- **cmd/server**: Точка входа и инициализация приложения
- **internal/config**: Управление конфигурацией
- **internal/handlers**: HTTP обработчики
- **internal/middleware**: HTTP middleware
- **internal/metrics**: Сбор и экспорт метрик
- **internal/models**: Модели данных
- **internal/server**: Настройка и управление HTTP сервером

### 2. Dependency Injection

Все компоненты получают свои зависимости через конструкторы:

```go
// Создание зависимостей
cfg := config.Load()
metrics := metrics.New()
handlers := handlers.New(cfg, metrics, &requestCount)
server := server.New(cfg)
```

### 3. Интерфейсы

Middleware реализуют общий интерфейс для лучшей композиции:

```go
type Middleware interface {
    Handler(next http.Handler) http.Handler
}
```

### 4. Конфигурация

Структурированная конфигурация с валидацией:

```go
type Config struct {
    Server   ServerConfig
    App      AppConfig
    Metrics  MetricsConfig
    Logging  LoggingConfig
}
```

## Улучшения после рефакторинга

### 1. Модульность
- Каждый пакет имеет четко определенную ответственность
- Легко добавлять новые функции без изменения существующего кода
- Возможность переиспользования компонентов

### 2. Тестируемость
- Unit тесты для каждого пакета
- Integration тесты для всего приложения
- Мокирование зависимостей через интерфейсы

### 3. Конфигурация
- Структурированная конфигурация с валидацией
- Поддержка различных окружений
- Типобезопасность

### 4. Метрики
- Инкапсулированные Prometheus метрики
- Легко добавлять новые метрики
- Опциональное включение/выключение

### 5. Middleware
- Композиция через интерфейсы
- Легко добавлять новые middleware
- Четкий порядок выполнения

### 6. Обработка ошибок
- Структурированная обработка ошибок
- Логирование с контекстом
- Graceful shutdown

## Запуск и тестирование

### Сборка
```bash
make build
```

### Тестирование
```bash
# Unit тесты
make test

# Тесты с покрытием
make coverage

# Integration тесты
go test ./test/integration/...
```

### Запуск
```bash
# Локально
./main

# Docker
make docker-build
make docker-run

# Docker Compose с мониторингом
make up
```

## Переменные окружения

| Переменная | Описание | По умолчанию |
|------------|----------|--------------|
| PORT | Порт сервера | 8080 |
| ENVIRONMENT | Окружение | development |
| APP_VERSION | Версия приложения | 1.0.0 |
| LOG_LEVEL | Уровень логирования | info |
| LOG_FORMAT | Формат логов | json |
| METRICS_ENABLED | Включить метрики | true |
| METRICS_PATH | Путь к Prometheus метрикам | /prometheus |
| READ_TIMEOUT | Таймаут чтения | 15s |
| WRITE_TIMEOUT | Таймаут записи | 15s |
| IDLE_TIMEOUT | Таймаут простоя | 60s |

## Endpoints

| Endpoint | Метод | Описание |
|----------|-------|----------|
| / | GET | Информация о сервере |
| /health | GET | Health check |
| /metrics | GET | Метрики в JSON формате |
| /prometheus | GET | Prometheus метрики |

## Мониторинг

Проект включает полный стек мониторинга:

- **Prometheus**: Сбор метрик
- **Grafana**: Визуализация
- **Alertmanager**: Уведомления
- **Node Exporter**: Системные метрики
- **Blackbox Exporter**: Проверка доступности

Доступ:
- Grafana: http://localhost:3000 (admin/admin123)
- Prometheus: http://localhost:9090
- Alertmanager: http://localhost:9093

## Безопасность

Добавлены security headers:
- X-Content-Type-Options: nosniff
- X-Frame-Options: DENY
- X-XSS-Protection: 1; mode=block
- Referrer-Policy: strict-origin-when-cross-origin

## Производительность

- Оптимизированный Docker образ с multi-stage build
- Статическая компиляция без CGO
- Graceful shutdown
- Настраиваемые таймауты
- Эффективное логирование 