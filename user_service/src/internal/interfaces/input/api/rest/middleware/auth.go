package middleware

import (
	"context"
	"user_service/src/pkg/utilities"
	"net/http"
)

func Authenticate(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		cookie,err:= r.Cookie("at")

		if err!=nil{
			http.Error(w,"Missing authorization token", http.StatusUnauthorized)
			return
		}

		claims, err:= utilities.ValidateJWT(cookie.Value)
		if err!=nil{
			http.Error(w, "Invalid token",http.StatusUnauthorized)
			return
		}

		ctx:= context.WithValue(r.Context(), "user", claims.Uid)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}