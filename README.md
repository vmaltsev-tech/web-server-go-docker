# 🚀 Go Web Server with Docker & Monitoring

Современный веб-сервер на Go с полным стеком мониторинга, следующий лучшим практикам 2025 года.

## ✨ Особенности

- ⚡ **Быстрый и эффективный** - Go 1.24.4 с оптимизированной сборкой
- 🐳 **Multi-stage Docker** - минимальный размер образа (Alpine Linux)
- 📊 **Полный мониторинг** - Prometheus + Grafana
- 🔒 **Безопасность** - security headers, non-root user, health checks
- 🧪 **Высокое покрытие тестами** - 1300+ строк тестов
- 🔄 **Graceful shutdown** - корректное завершение работы
- 📈 **Observability** - структурированные логи и метрики

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Web App    │────│   Prometheus    │────│    Grafana      │
│   Port: 8080    │    │   Port: 9090    │    │   Port: 3000    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │
         └─── Health Checks, Metrics, Security Headers
```

## 🚀 Быстрый старт

### Предварительные требования

- Docker & Docker Compose
- Go 1.24+ (для local development)
- Make (опционально)

### Запуск полного мониторинг стека

```bash
# Единый Docker Compose - запускает все сервисы
docker-compose up -d

# Проверка статуса всех сервисов
docker-compose ps

# Просмотр логов конкретных сервисов
docker-compose logs -f web          # Веб приложение
docker-compose logs -f prometheus   # Метрики
docker-compose logs -f grafana      # Дашборды
```

### Все сервисы в одном файле

Теперь все компоненты мониторинга объединены в единый `docker-compose.yaml`:
- 🚀 **Web Server** (порт 8080) - ваше Go приложение
- 🔍 **Prometheus** (порт 9090) - сбор метрик и алерты
- 📊 **Grafana** (порт 3000) - визуализация и дашборды
- 📈 **Node Exporter** (порт 9100) - системные метрики
- 🔍 **Blackbox Exporter** (порт 9115) - мониторинг доступности
- 🚨 **Alertmanager** (порт 9093) - уведомления

### Проверка работы

```bash
# Проверка здоровья приложения
curl http://localhost:8080/health

# Просмотр метрик
curl http://localhost:8080/metrics

# Если установлен Make - удобные команды
make monitoring-status  # Полный статус + алерты
make monitoring-check   # Быстрая проверка всех endpoint'ов
make monitoring-logs    # Логи мониторинга
make monitoring-reload  # Перезагрузка конфигурации Prometheus
make health            # Проверка здоровья
make metrics           # Просмотр метрик
```

### Доступ к сервисам

После запуска `docker-compose up -d` все сервисы будут доступны:

| Сервис | URL | Описание |
|--------|-----|----------|
| **Web Server** | http://localhost:8080 | Основное приложение |
| **Prometheus** | http://localhost:9090 | Сбор метрик и алерты |
| **Grafana** | http://localhost:3000 | Дашборды (admin/admin123) |
| **Node Exporter** | http://localhost:9100 | Системные метрики |
| **Blackbox Exporter** | http://localhost:9115 | Мониторинг доступности |
| **Alertmanager** | http://localhost:9093 | Управление уведомлениями |

## 📋 API Endpoints

| Endpoint | Method | Описание |
|----------|--------|----------|
| `/` | GET | Информация о сервере |
| `/health` | GET | Health check |
| `/metrics` | GET | Метрики приложения (JSON) |
| `/prometheus` | GET | Prometheus метрики |

### Примеры ответов

**GET /**
```json
{
  "message": "DevOps Portfolio 2025 - Go Web Server",
  "environment": "development",
  "port": "8080"
}
```

**GET /health**
```json
{
  "status": "OK",
  "timestamp": "2025-01-XX...",
  "version": "1.0.0"
}
```

## 🔧 Конфигурация

### Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|--------------|-----------|
| `PORT` | `8080` | Порт сервера |
| `ENVIRONMENT` | `development` | Окружение (development/staging/production/test) |
| `APP_VERSION` | `1.0.0` | Версия приложения |
| `LOG_LEVEL` | `info` | Уровень логирования |
| `READ_TIMEOUT` | `15s` | Read timeout (production) |
| `WRITE_TIMEOUT` | `15s` | Write timeout (production) |
| `IDLE_TIMEOUT` | `60s` | Idle timeout (production) |

### Production конфигурация

```bash
export ENVIRONMENT=production
export LOG_LEVEL=warn
export READ_TIMEOUT=30s
export WRITE_TIMEOUT=30s
export IDLE_TIMEOUT=120s
```

## 🧪 Тестирование

```bash
# Запуск всех тестов
make test

# Тесты с покрытием
make coverage

# Benchmark тесты
make benchmark

# Проверка безопасности
make security-scan

# Полная CI проверка
make ci
```

## 📊 Мониторинг

### Prometheus (http://localhost:9090)

**Метрики и алерты:**
- `http_requests_total` - общее количество HTTP запросов
- `http_request_duration_seconds` - время выполнения запросов  
- `server_uptime_seconds` - время работы сервера
- `go_memstats_*` - метрики памяти Go
- `go_goroutines` - количество горутин

**Активные алерты:**
- 🚨 WebServerDown - сервер недоступен
- ⚡ HighRequestLatency - высокая задержка запросов
- 🔥 HighErrorRate - высокий процент ошибок
- 💾 HighMemoryUsage - высокое потребление памяти
- 🔄 TooManyGoroutines - слишком много горутин
- ⏱️ HighResponseTime - медленное время ответа
- ❌ PrometheusTargetDown - цель мониторинга недоступна

### Grafana (http://localhost:3000)

- **Логин**: admin
- **Пароль**: admin123

Предустановленные дашборды для мониторинга производительности.

### Доступ к мониторингу

```bash
# Открыть Prometheus
make prometheus

# Открыть Grafana  
make grafana
```

## 🐳 Docker

### Сборка

```bash
# Development сборка
make docker-build

# Production сборка
make prod-build
```

### Особенности Docker образа

- **Multi-stage build** - уменьшение размера
- **Alpine Linux** - минимальная база
- **Static binary** - без внешних зависимостей
- **Non-root user** - безопасность
- **Health checks** - мониторинг состояния

## 🔒 Безопасность

### Реализованные меры

- ✅ Non-root пользователь в контейнере
- ✅ Security headers (X-Content-Type-Options, X-Frame-Options, X-XSS-Protection)
- ✅ HTTP method validation
- ✅ Proper error handling без утечки информации
- ✅ Graceful shutdown
- ✅ Input validation

### Security сканирование

```bash
make security-scan
```

## 🛠️ Разработка

### Локальная разработка

```bash
# Установка зависимостей
go mod download

# Запуск приложения
go run main.go config.go

# Или использование Make
make build
./main
```

### Code quality

```bash
# Форматирование кода
make fmt

# Проверка кода
make vet

# Полная проверка для разработки
make dev
```

## 📈 Производительность

### Benchmark результаты

```bash
make benchmark
```

Типичные результаты:
- Health endpoint: ~100k requests/sec
- Info endpoint: ~95k requests/sec
- Metrics endpoint: ~90k requests/sec

### Оптимизации

- Статическая компиляция с `-ldflags="-w -s"`
- Отключение CGO для лучшей производительности
- Proper connection pooling и timeouts
- Efficient JSON serialization

## 🚀 Deployment

### Kubernetes (рекомендуется для production)

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: web-server-go
spec:
  replicas: 3
  selector:
    matchLabels:
      app: web-server-go
  template:
    metadata:
      labels:
        app: web-server-go
    spec:
      containers:
      - name: web-server-go
        image: web-server-go:latest
        ports:
        - containerPort: 8080
        env:
        - name: ENVIRONMENT
          value: "production"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
```

### Docker Swarm

```bash
docker stack deploy -c docker-compose.yaml web-server-stack
```

## 📝 Логирование

Приложение использует структурированное логирование с поддержкой различных уровней:

- `DEBUG` - детальная информация для отладки
- `INFO` - общая информация о работе
- `WARN` - предупреждения
- `ERROR` - ошибки
- `FATAL` - критические ошибки

## 🤝 Contributing

1. Fork репозиторий
2. Создайте feature branch (`git checkout -b feature/amazing-feature`)
3. Commit изменения (`git commit -m 'Add amazing feature'`)
4. Push в branch (`git push origin feature/amazing-feature`)
5. Создайте Pull Request

### Code Style

- Используйте `gofmt` для форматирования
- Следуйте [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Добавляйте тесты для нового функционала
- Обновляйте документацию

## 📄 Лицензия

MIT License - см. [LICENSE](LICENSE) файл для деталей.

## 🔗 Полезные ссылки

- [Go Documentation](https://golang.org/doc/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)

---

**Версия**: 1.0.0  
**Go версия**: 1.24.4  
**Последнее обновление**: 2025 