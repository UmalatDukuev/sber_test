package handlers

import (
	"encoding/json"
	"net/http"

	"sber_test/internal/service"
)

// Execute handles the loan calculation request and returns the response.
func Execute(svc *service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.ExecuteRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, `{"error":"invalid request"}`, http.StatusBadRequest)
			return
		}

		resp, _, err := svc.Execute(req)
		if err != nil {
			http.Error(w, `{"error":"`+err.Error()+`"}`, http.StatusBadRequest)
			return
		}

		out := struct {
			Result service.ExecuteResponse `json:"result"`
		}{
			Result: resp,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(out); err != nil {
			http.Error(w, `{"error":"failed to encode data"}`, http.StatusInternalServerError)
			return
		}
	}
}
