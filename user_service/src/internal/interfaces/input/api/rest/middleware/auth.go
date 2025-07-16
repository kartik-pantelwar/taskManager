package middleware

import (
	"context"
	"net/http"
	errorhandling "user_service/src/pkg/error_handling"
	"user_service/src/pkg/utilities"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("at")

		if err != nil {
			errorhandling.HandleError(w, "Missing Authorization Token", http.StatusUnauthorized)
			return
		}

		claims, err := utilities.ValidateJWT(cookie.Value)
		if err != nil {
			errorhandling.HandleError(w, "Invalid Authorization Token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", claims.Uid)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
