# 📊 Monitoring Guide

Руководство по использованию единого мониторинг стека.

## 🚀 Быстрый старт

```bash
# Запуск всего стека
docker-compose up -d

# Проверка статуса
docker-compose ps

# Быстрая проверка всех сервисов
curl http://localhost:8080/health   # Web Server
curl http://localhost:9090/-/ready  # Prometheus
curl http://localhost:3000/api/health # Grafana
```

## 📈 Доступ к сервисам

| Сервис | URL | Логин | Назначение |
|--------|-----|-------|------------|
| **Web Server** | http://localhost:8080 | - | Основное приложение |
| **Prometheus** | http://localhost:9090 | - | Метрики и алерты |
| **Grafana** | http://localhost:3000 | admin/admin123 | Дашборды |
| **Alertmanager** | http://localhost:9093 | - | Уведомления |
| **Node Exporter** | http://localhost:9100 | - | Системные метрики |
| **Blackbox Exporter** | http://localhost:9115 | - | Проверка доступности |

## 🎯 Ключевые метрики

### Приложение
- `http_requests_total{method, status}` - количество HTTP запросов
- `http_request_duration_seconds` - время выполнения запросов
- `server_uptime_seconds` - время работы сервера

### Система (Node Exporter)
- `node_cpu_seconds_total` - загрузка CPU
- `node_memory_MemAvailable_bytes` - доступная память
- `node_filesystem_avail_bytes` - свободное место на диске
- `node_load1` - средняя загрузка системы

### Go Runtime
- `go_memstats_heap_alloc_bytes` - использование памяти heap
- `go_goroutines` - количество горутин
- `go_gc_duration_seconds` - время сборки мусора

## 🚨 Настроенные алерты

### Критические
- **WebServerDown** - сервер недоступен (>1 мин)
- **PrometheusTargetDown** - цель мониторинга недоступна (>1 мин)

### Предупреждения  
- **HighRequestLatency** - 95% запросов >0.5с (>2 мин)
- **HighErrorRate** - >10% ошибок 5xx (>2 мин)
- **HighMemoryUsage** - >100MB heap (>5 мин)
- **TooManyGoroutines** - >1000 горутин (>5 мин)
- **HighResponseTime** - 99% запросов >1с (>2 мин)

## 🔧 Полезные команды

### Docker Compose
```bash
# Управление стеком
docker-compose up -d           # Запуск
docker-compose down            # Остановка
docker-compose restart web    # Перезапуск сервиса
docker-compose logs -f web    # Логи сервиса

# Проверка статуса
docker-compose ps             # Статус всех сервисов
docker-compose top            # Процессы в контейнерах
```

### Prometheus API
```bash
# Проверка целей
curl http://localhost:9090/api/v1/targets

# Активные алерты
curl http://localhost:9090/api/v1/alerts

# Перезагрузка конфигурации
curl -X POST http://localhost:9090/-/reload

# Запрос метрик
curl "http://localhost:9090/api/v1/query?query=up"
```

### Grafana
```bash
# Проверка здоровья
curl http://localhost:3000/api/health

# Список дашбордов (требует аутентификации)
curl -u admin:admin123 http://localhost:3000/api/dashboards/home
```

## 📊 Создание дашбордов в Grafana

1. Откройте http://localhost:3000
2. Логин: `admin`, пароль: `admin123`
3. Перейдите в **Dashboards** → **New** → **New Dashboard**
4. Добавьте панель с источником данных **Prometheus**

### Популярные запросы для дашбордов

```promql
# RPS (requests per second)
rate(http_requests_total[5m])

# Время ответа (95-й процентиль)
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Использование CPU
rate(node_cpu_seconds_total{mode!="idle"}[5m]) * 100

# Использование памяти
(1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)) * 100

# Количество горутин
go_goroutines

# Размер heap памяти
go_memstats_heap_alloc_bytes / 1024 / 1024
```

## 🔍 Troubleshooting

### Сервис не запускается
```bash
# Проверка логов
docker-compose logs service-name

# Проверка конфигурации
docker-compose config
```

### Prometheus не видит цели
```bash
# Проверка сетевого подключения
docker-compose exec prometheus wget -O- http://web:8080/prometheus

# Проверка конфигурации
docker-compose exec prometheus cat /etc/prometheus/prometheus.yml
```

### Grafana недоступна
```bash
# Проверка статуса
docker-compose logs grafana

# Сброс пароля (если забыли)
docker-compose exec grafana grafana-cli admin reset-admin-password newpassword
```

## 📚 Дополнительные ресурсы

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Node Exporter Metrics](https://github.com/prometheus/node_exporter)
- [Go Application Metrics](https://prometheus.io/docs/guides/go-application/)

---

**Последнее обновление**: 2025-06-06
