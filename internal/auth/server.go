package auth

import (
	"context"
	"database/sql"
	"log/slog"

	"connectrpc.com/connect"
	v1 "github.com/damejeras/goose/api/gen/go/v1"
	"github.com/damejeras/goose/db/sqlc"
)

// Server implements the AuthService
type Server struct {
	authService *Service
	queries     *sqlc.Queries
	logger      *slog.Logger
}

// NewServer creates a new auth server
func NewServer(authService *Service, queries *sqlc.Queries, logger *slog.Logger) *Server {
	return &Server{
		authService: authService,
		queries:     queries,
		logger:      logger,
	}
}

// Login handles user login with Google ID token
func (s *Server) Login(ctx context.Context, req *connect.Request[v1.LoginRequest]) (*connect.Response[v1.LoginResponse], error) {
	// Validate Google ID token
	tokenInfo, err := s.authService.ValidateGoogleIDToken(ctx, req.Msg.GoogleIdToken)
	if err != nil {
		s.logger.Error("failed to validate Google ID token", "error", err)
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	if !tokenInfo.Verified {
		s.logger.Warn("unverified email attempted login", "email", tokenInfo.Email)
		return nil, connect.NewError(connect.CodePermissionDenied, ErrUnauthorized)
	}

	// Find or create user
	user, err := s.queries.FindUserByGoogleID(ctx, sql.NullString{
		String: tokenInfo.GoogleID,
		Valid:  true,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			// Create new user
			user, err = s.queries.CreateUser(ctx, sqlc.CreateUserParams{
				Email: tokenInfo.Email,
				GoogleID: sql.NullString{
					String: tokenInfo.GoogleID,
					Valid:  true,
				},
				Name: tokenInfo.Name,
			})
			if err != nil {
				s.logger.Error("failed to create user", "error", err)
				return nil, connect.NewError(connect.CodeInternal, err)
			}
			s.logger.Info("new user created", "user_id", user.ID, "email", user.Email)
		} else {
			s.logger.Error("failed to find user", "error", err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	} else {
		// Update existing user's profile info and last login
		if err := s.queries.UpdateUserProfile(ctx, sqlc.UpdateUserProfileParams{
			Name: tokenInfo.Name,
			ID:   user.ID,
		}); err != nil {
			s.logger.Warn("failed to update user profile", "user_id", user.ID, "error", err)
		}
		// Refresh user data
		user, err = s.queries.GetUser(ctx, user.ID)
		if err != nil {
			s.logger.Error("failed to get user after update", "error", err)
			return nil, connect.NewError(connect.CodeInternal, err)
		}
	}

	// Generate JWT
	jwt, err := s.authService.GenerateJWT(user.ID, user.Email)
	if err != nil {
		s.logger.Error("failed to generate JWT", "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	s.logger.Info("user logged in", "user_id", user.ID, "email", user.Email)

	return connect.NewResponse(&v1.LoginResponse{
		Jwt: jwt,
		User: &v1.User{
			Id:       user.ID,
			Email:    user.Email,
			GoogleId: user.GoogleID.String,
			Name:     user.Name,
		},
	}), nil
}

// GetCurrentUser returns the current authenticated user
func (s *Server) GetCurrentUser(ctx context.Context, req *connect.Request[v1.GetCurrentUserRequest]) (*connect.Response[v1.GetCurrentUserResponse], error) {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, ErrUnauthorized)
	}

	user, err := s.queries.GetUser(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, err)
		}
		s.logger.Error("failed to get user", "user_id", userID, "error", err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.GetCurrentUserResponse{
		User: &v1.User{
			Id:       user.ID,
			Email:    user.Email,
			GoogleId: user.GoogleID.String,
			Name:     user.Name,
		},
	}), nil
}

// Logout handles user logout (currently just returns success)
func (s *Server) Logout(ctx context.Context, req *connect.Request[v1.LogoutRequest]) (*connect.Response[v1.LogoutResponse], error) {
	// With JWT, logout is handled client-side by removing the token
	// Optionally, you could implement token blacklisting here
	userID, ok := GetUserIDFromContext(ctx)
	if ok {
		s.logger.Info("user logged out", "user_id", userID)
	}

	return connect.NewResponse(&v1.LogoutResponse{
		Success: true,
	}), nil
}
