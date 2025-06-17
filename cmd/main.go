package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/go-chi/chi"

	"sber_test/internal/handlers"
	"sber_test/internal/repo/cache"
	"sber_test/internal/service"
)

type Config struct {
	Port int `yaml:"port"`
}

func loadConfig(path string) Config {
	b, err := os.ReadFile(path)
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

	memCache := cache.New()

	svc := service.New(memCache)

	r := chi.NewRouter()
	handlers.RegisterRoutes(r, svc)

	addr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
