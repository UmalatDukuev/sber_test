package handlers

import (
	"encoding/json"
	"net/http"

	"sber_test/internal/service"
)

func GetCache(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		all := svc.GetAll()
		if len(all) == 0 {
			http.Error(w, `{"error":"empty cache"}`, http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(all)
	}
}
