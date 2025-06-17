package handlers

import (
	"sber_test/internal/service"

	"github.com/go-chi/chi"
)

func RegisterRoutes(r chi.Router, svc *service.Service) {
	r.Use(Logger)
	r.Post("/execute", Execute(svc))
	r.Get("/cache", GetCache(svc))
}
