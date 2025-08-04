package middleware

import (
	"net/http"
	"os"
)

func AuthKey(next http.Handler) http.Handler {
	secretKey := os.Getenv("SECRET_KEY")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("x-api-secret") != secretKey {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
