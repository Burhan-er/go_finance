package middleware

import (
	"encoding/json"
	"net/http"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationMiddleware struct{}

func NewValidationMiddleware() *ValidationMiddleware {
	return &ValidationMiddleware{}
}

func (m *ValidationMiddleware) ValidateJSON(v interface{}) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Body == nil {
				http.Error(w, "Request body is empty", http.StatusBadRequest)
				return
			}

			err := json.NewDecoder(r.Body).Decode(v)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Invalid JSON format: " + err.Error(),
				})
				return
			}

			r.Body = http.MaxBytesReader(w, r.Body, 1048576) // 1MB limit

			next.ServeHTTP(w, r)
		})
	}
}

func (m *ValidationMiddleware) ValidateQueryParams(required ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			query := r.URL.Query()
			var missingParams []string

			for _, param := range required {
				if query.Get(param) == "" {
					missingParams = append(missingParams, param)
				}
			}

			if len(missingParams) > 0 {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"error":   "Missing required query parameters",
					"missing": missingParams,
				})
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
