package apikey

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"connectrpc.com/connect"
	v1 "github.com/damejeras/goose/api/gen/go/v1"
	"github.com/damejeras/goose/db/sqlc"
	"github.com/damejeras/goose/internal/auth"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Server implements the APIKeyService
type Server struct {
	queries *sqlc.Queries
	logger  *slog.Logger
}

// NewServer creates a new API key server
func NewServer(queries *sqlc.Queries, logger *slog.Logger) *Server {
	return &Server{
		queries: queries,
		logger:  logger,
	}
}

// CreateAPIKey generates a new API key for the user
func (s *Server) CreateAPIKey(ctx context.Context, req *connect.Request[v1.CreateAPIKeyRequest]) (*connect.Response[v1.CreateAPIKeyResponse], error) {
	// Get user ID from context (set by auth middleware)
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not authenticated"))
	}

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	// generateKey new API key
	id, key, err := generateKey()
	if err != nil {
		s.logger.Error("failed to generate API key", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to generate API key"))
	}

	// hash the key for storage
	keyHash := hash(key)

	// Extract prefix and suffix
	prefix, suffix := ExtractParts(key)

	// Store in database
	dbKey, err := s.queries.CreateAPIKey(ctx, sqlc.CreateAPIKeyParams{
		ID:        id,
		UserID:    userID,
		Name:      req.Msg.Name,
		KeyHash:   keyHash,
		KeyPrefix: prefix,
		KeySuffix: suffix,
	})
	if err != nil {
		s.logger.Error("failed to create API key", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create API key"))
	}

	return connect.NewResponse(&v1.CreateAPIKeyResponse{
		Id:        dbKey.ID,
		Name:      dbKey.Name,
		Key:       key, // Return the full key - only time it's shown
		CreatedAt: timestamppb.New(dbKey.CreatedAt),
	}), nil
}

// ListAPIKeys returns all API keys for the authenticated user
func (s *Server) ListAPIKeys(ctx context.Context, req *connect.Request[v1.ListAPIKeysRequest]) (*connect.Response[v1.ListAPIKeysResponse], error) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not authenticated"))
	}

	// Get API keys from database
	dbKeys, err := s.queries.ListAPIKeysByUserID(ctx, userID)
	if err != nil {
		s.logger.Error("failed to list API keys", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list API keys"))
	}

	// Convert to proto messages with masked keys
	apiKeys := make([]*v1.APIKey, len(dbKeys))
	for i, dbKey := range dbKeys {
		// Reconstruct the masked key from prefix and suffix
		maskedKey := fmt.Sprintf("%s****...****%s", dbKey.KeyPrefix, dbKey.KeySuffix)

		var lastUsedAt *timestamppb.Timestamp
		if dbKey.LastUsedAt.Valid {
			lastUsedAt = timestamppb.New(dbKey.LastUsedAt.Time)
		}

		apiKeys[i] = &v1.APIKey{
			Id:         dbKey.ID,
			Name:       dbKey.Name,
			KeyMasked:  maskedKey,
			CreatedAt:  timestamppb.New(dbKey.CreatedAt),
			LastUsedAt: lastUsedAt,
		}
	}

	return connect.NewResponse(&v1.ListAPIKeysResponse{
		ApiKeys: apiKeys,
	}), nil
}

// DeleteAPIKey deletes an API key
func (s *Server) DeleteAPIKey(ctx context.Context, req *connect.Request[v1.DeleteAPIKeyRequest]) (*connect.Response[v1.DeleteAPIKeyResponse], error) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not authenticated"))
	}

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("id is required"))
	}

	// Delete the API key
	err := s.queries.DeleteAPIKey(ctx, sqlc.DeleteAPIKeyParams{
		ID:     req.Msg.Id,
		UserID: userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("API key not found"))
		}
		s.logger.Error("failed to delete API key", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete API key"))
	}

	return connect.NewResponse(&v1.DeleteAPIKeyResponse{
		Success: true,
	}), nil
}

// UpdateAPIKey updates an API key (name only)
func (s *Server) UpdateAPIKey(ctx context.Context, req *connect.Request[v1.UpdateAPIKeyRequest]) (*connect.Response[v1.UpdateAPIKeyResponse], error) {
	// Get user ID from context
	userID, ok := auth.GetUserIDFromContext(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("user not authenticated"))
	}

	if req.Msg.Id == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("id is required"))
	}

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name is required"))
	}

	// Update the API key
	dbKey, err := s.queries.UpdateAPIKeyName(ctx, sqlc.UpdateAPIKeyNameParams{
		Name:   req.Msg.Name,
		ID:     req.Msg.Id,
		UserID: userID,
	})
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("API key not found"))
		}
		s.logger.Error("failed to update API key", "error", err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update API key"))
	}

	// Reconstruct the masked key
	maskedKey := fmt.Sprintf("%s****...****%s", dbKey.KeyPrefix, dbKey.KeySuffix)

	var lastUsedAt *timestamppb.Timestamp
	if dbKey.LastUsedAt.Valid {
		lastUsedAt = timestamppb.New(dbKey.LastUsedAt.Time)
	}

	return connect.NewResponse(&v1.UpdateAPIKeyResponse{
		ApiKey: &v1.APIKey{
			Id:         dbKey.ID,
			Name:       dbKey.Name,
			KeyMasked:  maskedKey,
			CreatedAt:  timestamppb.New(dbKey.CreatedAt),
			LastUsedAt: lastUsedAt,
		},
	}), nil
}
