package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"sber_test/internal/repo/cache"
	"sber_test/internal/service"
)

func TestHandlers(t *testing.T) {
	cacheService := cache.New()
	svc := service.New(cacheService)
	err := os.Chdir("../../")
	if err != nil {
		panic("failed to change directory: " + err.Error())
	}
	t.Run("Test GetCache", func(t *testing.T) {
		handler := GetCache(svc)
		cacheService.Add(service.CacheItem{ID: 1})

		req, err := http.NewRequest("GET", "/cache", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
		}

		cacheService = cache.New()
		svc = service.New(cacheService)
		handler = GetCache(svc)
		req, err = http.NewRequest("GET", "/cache", nil)
		if err != nil {
			t.Fatal(err)
		}
		rr = httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("expected status %v, got %v", http.StatusBadRequest, rr.Code)
		}
	})

	t.Run("Test Execute", func(t *testing.T) {
		handler := Execute(svc)

		tests := []struct {
			body       service.ExecuteRequest
			statusCode int
		}{
			{
				body: service.ExecuteRequest{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
					Program: map[string]bool{
						"military": true,
					},
				},
				statusCode: http.StatusOK,
			},
			{
				body: service.ExecuteRequest{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
					Program: map[string]bool{
						"invalid": true,
					},
				},
				statusCode: http.StatusBadRequest,
			},
			{
				body: service.ExecuteRequest{
					ObjectCost:     5000000,
					InitialPayment: 1000000,
					Months:         240,
					Program:        map[string]bool{},
				},
				statusCode: http.StatusBadRequest,
			},
			{
				body: service.ExecuteRequest{
					ObjectCost:     5000000,
					InitialPayment: 50000,
					Months:         240,
					Program: map[string]bool{
						"military": true,
					},
				},
				statusCode: http.StatusBadRequest,
			},
		}

		for _, tt := range tests {
			bodyBytes, err := json.Marshal(tt.body)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest("POST", "/execute", bytes.NewBuffer(bodyBytes))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code != tt.statusCode {
				t.Errorf("expected status %v, got %v", tt.statusCode, rr.Code)
			}
		}
	})

	t.Run("Test Logger", func(t *testing.T) {
		handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))

		req, err := http.NewRequest("GET", "/cache", nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("expected status %v, got %v", http.StatusOK, rr.Code)
		}
	})
}
