package middleware

import (
	"context"
	"net/http"
	"strconv"
	pb "task_service/src/internal/interfaces/input/grpc/generated/generated"
	errorhandling "task_service/src/pkg/error_handling"
)

func SessionAuthMiddleware(grpcClient pb.SessionValidatorClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract "sess" cookie (refresh token/session ID)
			cookie, err := r.Cookie("sess")
			if err != nil {
				errorhandling.HandleError(w, "Session Cookie is missing", http.StatusUnauthorized)
				return
			}
			sessionID := cookie.Value
			// gRPC call to user service for session validation
			resp, err := grpcClient.ValidateSession(context.Background(), &pb.ValidateSessionRequest{
				SessionId: sessionID,
			})
			if err != nil || !resp.Valid {
				errorhandling.HandleError(w, "invalid session", http.StatusUnauthorized)
				return
			}

			// Converting user id into int
			userID, err := strconv.Atoi(resp.UserId)
			if err != nil {
				errorhandling.HandleError(w, "invalid user ID from session", http.StatusInternalServerError)
				return
			}
			ctx := context.WithValue(r.Context(), "user_id", userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
