package main

import (
	"log"

	"web-server-go-docker/internal/config"
	"web-server-go-docker/internal/server"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Создаем сервер
	srv, err := server.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	// Запускаем сервер
	if err := srv.Start(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
} 