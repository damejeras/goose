package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"connectrpc.com/connect"
	"github.com/damejeras/goose/api/gen/go/v1/v1connect"
	"github.com/damejeras/goose/db"
	"github.com/damejeras/goose/db/sqlc"
	"github.com/damejeras/goose/frontend"
	"github.com/damejeras/goose/internal/apikey"
	"github.com/damejeras/goose/internal/auth"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// Parse command line flags
	googleClientID := flag.String("google-client-id", os.Getenv("GOOGLE_CLIENT_ID"), "Google OAuth client ID")
	jwtSecretStr := flag.String("jwt-secret", os.Getenv("JWT_SECRET"), "JWT secret (base64 encoded)")
	dbPath := flag.String("db", "storage/goose.db", "Database path")
	port := flag.String("port", "8080", "Server port")
	flag.Parse()

	// Setup logger
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// Validate required config
	if *googleClientID == "" {
		logger.Error("GOOGLE_CLIENT_ID is required")
		os.Exit(1)
	}

	var jwtSecret []byte
	var err error
	if *jwtSecretStr != "" {
		jwtSecret, err = base64.StdEncoding.DecodeString(*jwtSecretStr)
		if err != nil {
			logger.Error("failed to decode JWT secret", "error", err)
			os.Exit(1)
		}
	} else {
		// generateKey random secret for development
		jwtSecret, err = auth.GenerateRandomSecret()
		if err != nil {
			logger.Error("failed to generate JWT secret", "error", err)
			os.Exit(1)
		}
		logger.Warn("using randomly generated JWT secret - tokens will not persist across restarts")
	}

	// Open database
	database, err := db.Open(context.Background(), logger, *dbPath)
	if err != nil {
		logger.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	queries := sqlc.New(database)

	authService := auth.NewService(auth.Config{
		GoogleClientID: *googleClientID,
		JWTSecret:      jwtSecret,
		JWTExpiration:  24 * time.Hour,
	})

	// Create auth interceptor - specify public methods that don't require auth
	publicMethods := []string{
		"/api.v1.AuthService/Login",
		"/api.v1.GreeterService/SayHello", // Keep greeter public for testing
	}
	authInterceptor := auth.NewInterceptor(authService, publicMethods)

	// Setup HTTP mux
	mux := http.NewServeMux()

	// Register auth service with interceptor
	authPath, authHandler := v1connect.NewAuthServiceHandler(
		auth.NewServer(authService, queries, logger),
		connect.WithInterceptors(authInterceptor),
	)

	mux.Handle(authPath, authHandler)

	// Register API key service with interceptor (requires authentication)
	apiKeyPath, apiKeyHandler := v1connect.NewAPIKeyServiceHandler(
		apikey.NewServer(queries, logger),
		connect.WithInterceptors(authInterceptor),
	)
	mux.Handle(apiKeyPath, apiKeyHandler)

	mux.Handle("/", frontend.Handler())

	addr := ":" + *port
	logger.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
