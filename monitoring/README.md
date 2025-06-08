# 📊 Monitoring Infrastructure

Полная система мониторинга для Go веб-сервера, включающая Prometheus, Grafana, Alertmanager и экспортеры.

## 🏗️ Архитектура

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Go Web App    │────│   Prometheus    │────│    Grafana      │
│   Port: 8080    │    │   Port: 9090    │    │   Port: 3000    │
└─────────────────┘    └─────────────────┘    └─────────────────┘
         │                        │                        │
         │                        │                        │
    ┌─────────┐            ┌─────────────┐         ┌─────────────┐
    │ Health  │            │ Alertmanager│         │ Dashboards  │
    │ Checks  │            │ Port: 9093  │         │ & Panels    │
    └─────────┘            └─────────────┘         └─────────────┘
         │                        │
         │                        │
    ┌─────────────────────────────────────────┐
    │           Exporters                     │
    │  ┌─────────────┐  ┌─────────────────┐   │
    │  │Node Exporter│  │Blackbox Exporter│   │
    │  │Port: 9100   │  │Port: 9115       │   │
    │  └─────────────┘  └─────────────────┘   │
    └─────────────────────────────────────────┘
```

## 📁 Структура директорий

```
monitoring/
├── prometheus/
│   ├── config/
│   │   └── prometheus.yml          # Основная конфигурация Prometheus
│   ├── rules/
│   │   └── alerts.yml              # Правила алертов
│   └── data/                       # Данные Prometheus (volume mount)
├── grafana/
│   ├── config/
│   │   └── provisioning/           # Автоматическая настройка
│   │       ├── datasources/        # Источники данных
│   │       └── dashboards/         # Настройка дашбордов
│   ├── dashboards/                 # JSON файлы дашбордов
│   ├── plugins/                    # Плагины Grafana
│   └── data/                       # Данные Grafana (volume mount)
├── alertmanager/
│   ├── config/
│   │   └── alertmanager.yml        # Конфигурация уведомлений
│   └── data/                       # Данные Alertmanager
├── exporters/
│   ├── node/                       # Node Exporter конфигурация
│   └── blackbox/
│       └── blackbox.yml            # Blackbox Exporter конфигурация
├── docker/
│   └── docker-compose.monitoring.yml # Отдельный compose для мониторинга
├── scripts/
│   └── setup.sh                    # Скрипт инициализации
└── README.md                       # Эта документация
```

## 🚀 Быстрый старт

### 1. Инициализация

```bash
# Запуск скрипта настройки
./monitoring/scripts/setup.sh
```

### 2. Запуск мониторинга

```bash
# Использование основного docker-compose
docker-compose up -d

# Или использование отдельного файла мониторинга
cd monitoring/docker
docker-compose -f docker-compose.monitoring.yml up -d
```

### 3. Доступ к сервисам

- **Grafana**: http://localhost:3000 (admin/admin123)
- **Prometheus**: http://localhost:9090
- **Alertmanager**: http://localhost:9093
- **Node Exporter**: http://localhost:9100
- **Blackbox Exporter**: http://localhost:9115
- **Web Server**: http://localhost:8080

## 📊 Метрики

### Go Application Metrics
- `http_requests_total` - Общее количество HTTP запросов
- `http_request_duration_seconds` - Время выполнения запросов
- `server_uptime_seconds` - Время работы сервера
- `go_goroutines` - Количество горутин
- `go_memstats_*` - Метрики памяти Go

### System Metrics (Node Exporter)
- `node_cpu_seconds_total` - Использование CPU
- `node_memory_*` - Метрики памяти
- `node_filesystem_*` - Метрики файловой системы
- `node_network_*` - Сетевые метрики

### Availability Metrics (Blackbox)
- `probe_success` - Успешность проверки
- `probe_duration_seconds` - Время отклика
- `probe_http_status_code` - HTTP статус код

## 🔔 Алерты

### Настроенные алерты:

1. **WebServerDown** - Сервер недоступен (критический)
2. **HighRequestLatency** - Высокая задержка запросов
3. **HighErrorRate** - Высокий процент ошибок
4. **HighMemoryUsage** - Высокое использование памяти
5. **TooManyGoroutines** - Слишком много горутин

### Каналы уведомлений:
- Email
- Slack
- Webhook

## 📈 Дашборды

### Go Web Server Monitoring
- HTTP Request Rate
- Request Duration (95th/50th percentile)
- Active Goroutines
- Memory Usage (Heap)
- Server Uptime
- Garbage Collector Duration

### Node Exporter Dashboard
- CPU Usage
- Memory Usage
- Disk I/O
- Network Traffic
- Load Average

## 🔧 Конфигурация

### Добавление новых targets

Отредактируйте `prometheus/config/prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'my-new-service'
    static_configs:
      - targets: ['service-host:port']
    metrics_path: '/metrics'
    scrape_interval: 15s
```

### Создание новых алертов

Добавьте правила в `prometheus/rules/alerts.yml`:

```yaml
- alert: MyAlert
  expr: my_metric > threshold
  for: 2m
  labels:
    severity: warning
  annotations:
    summary: "Alert summary"
    description: "Alert description"
```

### Настройка уведомлений

Отредактируйте `alertmanager/config/alertmanager.yml`:

```yaml
receivers:
  - name: 'my-receiver'
    slack_configs:
      - api_url: 'YOUR_SLACK_WEBHOOK'
        channel: '#alerts'
```

## 🛠️ Команды управления

```bash
# Просмотр статуса сервисов
docker-compose ps

# Просмотр логов
docker-compose logs -f [service_name]

# Перезагрузка конфигурации Prometheus
curl -X POST http://localhost:9090/-/reload

# Перезагрузка конфигурации Alertmanager
curl -X POST http://localhost:9093/-/reload

# Остановка мониторинга
docker-compose down

# Остановка с удалением данных
docker-compose down -v
```

## 🔍 Troubleshooting

### Prometheus не может подключиться к targets

1. Проверьте network connectivity:
```bash
docker-compose exec prometheus wget -qO- http://web:8080/prometheus
```

2. Проверьте конфигурацию:
```bash
docker-compose exec prometheus promtool check config /etc/prometheus/prometheus.yml
```

### Grafana не показывает данные

1. Проверьте datasource:
   - Перейдите в Configuration → Data Sources
   - Проверьте URL: `http://prometheus:9090`
   - Нажмите "Test"

2. Проверьте доступность Prometheus:
```bash
curl http://localhost:9090/api/v1/targets
```

### Alertmanager не отправляет уведомления

1. Проверьте конфигурацию:
```bash
docker-compose exec alertmanager amtool check-config /etc/alertmanager/alertmanager.yml
```

2. Проверьте статус алертов:
```bash
curl http://localhost:9093/api/v1/alerts
```

## 📚 Дополнительные ресурсы

- [Prometheus Documentation](https://prometheus.io/docs/)
- [Grafana Documentation](https://grafana.com/docs/)
- [Alertmanager Documentation](https://prometheus.io/docs/alerting/latest/alertmanager/)
- [Go Application Monitoring Best Practices](https://prometheus.io/docs/guides/go-application/)

## 🤝 Поддержка

Для вопросов и проблем создайте issue в репозитории или обратитесь к команде DevOps.

---

**Версия**: 1.0.0  
**Последнее обновление**: 2025  
**Статус**: Production Ready ✅ 