
FROM golang:1.24.4-alpine3.22 AS builder

WORKDIR /app

# Копируем только go.mod (go.sum отсутствует в проекте)
COPY go.mod ./

# Загружаем зависимости (убираем go mod verify так как нет go.sum)
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем оптимизированный статический бинарник
RUN CGO_ENABLED=0 GOOS=linux go build \
    -a -installsuffix cgo \
    -ldflags="-w -s" \
    -o main ./cmd/server

# Финальный образ
FROM alpine:latest

# Устанавливаем ca-certificates и curl для healthcheck
RUN apk --no-cache add ca-certificates curl

# Исправляем рабочую директорию
WORKDIR /app

# Правильно копируем бинарник
COPY --from=builder /app/main ./main

# Создаем пользователя и назначаем права
RUN adduser -D -s /bin/sh appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 8080

# Исправляем healthcheck - используем curl вместо wget
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/health || exit 1

CMD ["./main"]
