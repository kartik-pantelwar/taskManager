package middleware

import (
	"context"
	"net/http"
	pb "task_service/src/internal/interfaces/input/grpc/generated/generated"
)

type SessionAuthMiddleware struct {
	GrpcClient pb.SessionValidatorClient
}

// SessionIDMiddleware extracts session_id from the request header and adds it to the request context.

func (m *SessionAuthMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionID := r.Header.Get("Session-Id")
		if sessionID == "" {
			http.Error(w, "Missing session ID", http.StatusUnauthorized)
			return
		}

		// *gRPC call to user service
		resp, err := m.GrpcClient.ValidateSession(context.Background(), &pb.ValidateSessionRequest{
			SessionId: sessionID,
		})
		if err != nil || !resp.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// *Optionally, add user_id to context
		ctx := context.WithValue(r.Context(), "user_id", resp.UserId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
