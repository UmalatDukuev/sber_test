// Package handlers defines the HTTP handlers for the application.
package handlers

import (
	"sber_test/internal/service"

	"github.com/go-chi/chi"
)

// RegisterRoutes registers HTTP routes for the application.
func RegisterRoutes(r chi.Router, svc *service.Service) {
	r.Use(Logger)
	r.Post("/execute", Execute(svc))
	r.Get("/cache", GetCache(svc))
}
