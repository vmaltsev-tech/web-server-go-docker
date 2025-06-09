# Устранение неполадок мониторинга

## Ошибки Node Exporter

### "connection reset by peer" или "broken pipe"

**Симптомы:**
```
node-exporter | time=2025-06-08T06:08:50.452Z level=ERROR source=http.go:219 msg="error encoding and sending metric family: write tcp 127.0.0.1:9100->127.0.0.1:46066: write: connection reset by peer"
```

**Причина:**
Эти ошибки возникают когда Prometheus разрывает соединение с node-exporter до завершения передачи всех метрик. Это может происходить по следующим причинам:

1. **Таймаут scrape** - Prometheus ждет ответ только определенное время
2. **Большой объем метрик** - node-exporter генерирует много системных метрик
3. **Сетевые проблемы** - временные проблемы с сетью в Docker

**Решения:**

#### 1. Увеличение таймаутов
```yaml
# prometheus.yml
- job_name: 'node-exporter'
  scrape_interval: 60s
  scrape_timeout: 55s
```

#### 2. Ограничение коллекторов
```yaml
# docker-compose.yml
command:
  - '--collector.disable-defaults'
  - '--collector.cpu'
  - '--collector.filesystem'
  - '--collector.loadavg'
  - '--collector.meminfo'
  - '--collector.time'
  - '--collector.uname'
```
Отключение `netdev` позволяет сократить объем сетевых метрик и снизить
нагрузку на обработку, что также помогает избегать прерываний соединения.

#### 3. Настройка логирования (уже применено)
```yaml
command:
  - '--log.level=warn'  # Снижает уровень логирования
```

**Важно:** Эти ошибки обычно не критичны и не влияют на сбор метрик. Prometheus повторяет запросы согласно расписанию.

### Проверка работоспособности

```bash
# Проверить доступность node-exporter
curl http://localhost:9100/metrics

# Проверить статус в Prometheus
curl http://localhost:9090/api/v1/targets

# Проверить количество метрик
curl -s http://localhost:9100/metrics | wc -l

# Проверить время ответа
time curl -s http://localhost:9100/metrics > /dev/null
```

### Мониторинг ошибок

Если ошибки происходят слишком часто (более 10% запросов), рассмотрите:

1. **Увеличение ресурсов** контейнера node-exporter
2. **Дальнейшее ограничение коллекторов**
3. **Увеличение scrape_interval** до 60s и `scrape_timeout` до 55s

## Другие проблемы

### Prometheus не может подключиться к targets

**Проверить:**
1. Сетевые настройки Docker
2. Правильность имен сервисов в docker-compose.yml
3. Порты и firewall

### Grafana не показывает данные

**Проверить:**
1. Настройки datasource в Grafana
2. Доступность Prometheus из Grafana
3. Правильность запросов в дашбордах

### Высокое потребление ресурсов

**Оптимизация:**
1. Ограничить retention в Prometheus
2. Уменьшить частоту scraping
3. Отключить ненужные коллекторы 