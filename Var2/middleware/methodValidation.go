package middleware

import (
	"net/http"
)

func MethodValidation(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet && r.Method != http.MethodPost && r.Method != http.MethodPut {
			http.Error(w, "Invalid HTTP method", http.StatusMethodNotAllowed)
			return
		}
		next(w, r)
	}
}
