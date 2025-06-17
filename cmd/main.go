// Package main starts the app.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sber_test/internal/handlers"
	"sber_test/internal/repo/cache"
	"sber_test/internal/service"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"gopkg.in/yaml.v2"
)

// Config struct.
type Config struct {
	Port int `yml:"port"`
}

// BasePath - safe path.
const BasePath = "C:\\GolangProgs\\sber_test"

func loadConfig(path string) Config {
	absPath := filepath.Join(BasePath, path)
	cleanPath := filepath.Clean(absPath)

	if !strings.HasPrefix(cleanPath, BasePath) {
		log.Fatalf("invalid path: %v", cleanPath)
	}

	b, err := os.ReadFile(cleanPath)
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}
	var cfg Config
	if err := yaml.Unmarshal(b, &cfg); err != nil {
		log.Fatalf("failed to parse config: %v", err)
	}
	return cfg
}

func main() {
	cfg := loadConfig("config.yml")
	addr := fmt.Sprintf(":%d", cfg.Port)

	// Создаём экземпляр кеша (предположительно, Cache уже настроен)
	c := cache.New()

	// Создаем экземпляр Service, передавая в него кеш
	svc := service.New(c)

	// Создаём новый маршрутизатор chi
	r := chi.NewRouter()

	// Регистрируем маршруты через handler
	handlers.RegisterRoutes(r, svc)

	// Настроить сервер
	server := &http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("Server is running on %s\n", addr)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
