package handlers

import (
	"encoding/json"
	"net/http"

	"sber_test/internal/service"
)

// GetCache returns the cached items in JSON format.
func GetCache(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		all := svc.GetAll()
		if len(all) == 0 {
			http.Error(w, `{"error":"empty cache"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(all); err != nil {
			http.Error(w, `{"error":"failed to encode data"}`, http.StatusInternalServerError)
			return
		}
	}
}
