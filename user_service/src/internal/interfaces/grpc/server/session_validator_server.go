package grpc

import (
	"context"
	"strconv"
	"time"
	"user_service/src/internal/adaptors/persistance"
	pb "user_service/src/internal/interfaces/grpc/generated/generated"
)

type SessionValidatorServer struct {
	pb.UnimplementedSessionValidatorServer
	sessionRepo persistance.SessionRepo
}

// NewSessionValidatorServer creates a new SessionValidatorServer instance
func NewSessionValidatorServer(sessionRepo persistance.SessionRepo) *SessionValidatorServer {
	return &SessionValidatorServer{
		sessionRepo: sessionRepo,
	}
}

func (s *SessionValidatorServer) ValidateSession(ctx context.Context, req *pb.ValidateSessionRequest) (*pb.ValidateSessionResponse, error) {
	sessionID := req.GetSessionId()

	// Check if session ID is empty
	if sessionID == "" {
		return &pb.ValidateSessionResponse{
			Valid: false,
			Error: "session_id is empty",
		}, nil
	}

	// Validate session ID by checking if it exists in database
	session, err := s.sessionRepo.GetSession(sessionID)
	if err != nil {
		return &pb.ValidateSessionResponse{
			Valid: false,
			Error: "session not found in database",
		}, nil
	}

	// Check if session is expired
	if time.Now().After(session.ExpiresAt) {
		return &pb.ValidateSessionResponse{
			Valid: false,
			Error: "session expired",
		}, nil
	}

	// Session is valid, return user ID
	return &pb.ValidateSessionResponse{
		Valid:  true,
		UserId: strconv.Itoa(session.Uid),
		Error:  "",
	}, nil
}

func (s *SessionValidatorServer) ValidateUser(ctx context.Context, req *pb.ValidateUserRequest) (*pb.ValidateUserResponse, error) {
	userIDStr := req.GetUserId()

	// Check if user ID is empty
	if userIDStr == "" {
		return &pb.ValidateUserResponse{
			Status: false,
		}, nil
	}

	// Convert user ID from string to int
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return &pb.ValidateUserResponse{
			Status: false,
		}, nil
	}

	// Check if user exists in database
	exists, err := s.sessionRepo.UserExists(userID)
	if err != nil {
		return &pb.ValidateUserResponse{
			Status: false,
		}, nil
	}

	return &pb.ValidateUserResponse{
		Status: exists,
	}, nil
}
