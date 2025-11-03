package auth

import (
	"context"
	"strings"

	"connectrpc.com/connect"
)

// Interceptor is a Connect RPC interceptor that validates JWT tokens
type Interceptor struct {
	authService *Service
	publicMethods map[string]bool
}

// NewInterceptor creates a new auth interceptor
func NewInterceptor(authService *Service, publicMethods []string) *Interceptor {
	publicMap := make(map[string]bool)
	for _, method := range publicMethods {
		publicMap[method] = true
	}
	return &Interceptor{
		authService:   authService,
		publicMethods: publicMap,
	}
}

type contextKey string

const UserIDContextKey contextKey = "user_id"

// WrapUnary wraps unary RPC calls with authentication
func (i *Interceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		// Check if this is a public method
		procedure := req.Spec().Procedure
		if i.publicMethods[procedure] {
			return next(ctx, req)
		}

		// Extract token from Authorization header
		auth := req.Header().Get("Authorization")
		if auth == "" {
			return nil, connect.NewError(connect.CodeUnauthenticated, ErrMissingToken)
		}

		// Remove "Bearer " prefix
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			return nil, connect.NewError(connect.CodeUnauthenticated, ErrInvalidToken)
		}

		// Validate JWT
		claims, err := i.authService.ValidateJWT(token)
		if err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}

		// Add user ID to context
		ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)

		return next(ctx, req)
	}
}

// WrapStreamingClient wraps streaming client calls with authentication
func (i *Interceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next // no-op for now, can add auth logic later if needed
}

// WrapStreamingHandler wraps streaming handler calls with authentication
func (i *Interceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next // no-op for now, can add auth logic later if needed
}

// GetUserIDFromContext extracts the user ID from the context
func GetUserIDFromContext(ctx context.Context) (int64, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(int64)
	return userID, ok
}
