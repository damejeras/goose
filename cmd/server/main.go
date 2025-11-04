package main

import (
	"context"
	"encoding/base64"
	"flag"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	devMode := flag.Bool("dev", false, "Enable development mode with Vite proxy")
	viteURL := flag.String("vite-url", "http://localhost:5173", "Vite dev server URL")
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

	// Setup frontend handler - proxy to Vite in dev mode, serve static files in production
	// Use "/{path...}" pattern to match all remaining requests (catch-all)
	if *devMode {
		logger.Info("development mode enabled", "vite_url", *viteURL)
		mux.Handle("/{path...}", newViteProxy(*viteURL, logger))
	} else {
		mux.Handle("/{path...}", frontend.Handler())
	}

	addr := ":" + *port
	logger.Info("server starting", "addr", addr)
	if err := http.ListenAndServe(addr, h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}

// newViteProxy creates a proxy handler that forwards requests to the Vite dev server.
// It handles both HTTP requests and WebSocket connections (for HMR).
func newViteProxy(target string, logger *slog.Logger) http.Handler {
	targetURL, err := url.Parse(target)
	if err != nil {
		logger.Error("failed to parse vite URL", "error", err)
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Invalid Vite URL", http.StatusInternalServerError)
		})
	}

	client := &http.Client{Timeout: 30 * time.Second}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Handle WebSocket upgrade requests for HMR
		if r.Header.Get("Upgrade") == "websocket" {
			proxyWebSocket(w, r, targetURL, logger)
			return
		}

		// Proxy HTTP request to Vite
		proxyURL := *targetURL
		proxyURL.Path = r.URL.Path
		proxyURL.RawQuery = r.URL.RawQuery

		proxyReq, _ := http.NewRequest(r.Method, proxyURL.String(), r.Body)
		proxyReq.Header = r.Header.Clone()

		resp, err := client.Do(proxyReq)
		if err != nil {
			logger.Error("vite dev server unavailable", "error", err)
			http.Error(w, "Vite dev server not available", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// If 404 and looks like a SPA route (no file extension), serve index.html
		if resp.StatusCode == http.StatusNotFound && !strings.Contains(r.URL.Path, ".") {
			proxyURL.Path = "/"
			indexReq, _ := http.NewRequest("GET", proxyURL.String(), nil)
			resp.Body.Close() // Close the 404 response

			resp, err = client.Do(indexReq)
			if err != nil {
				http.Error(w, "Failed to load index.html", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()
		}

		// Copy response to client
		for key, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(key, value)
			}
		}
		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})
}

// proxyWebSocket handles WebSocket connections for Vite HMR
func proxyWebSocket(w http.ResponseWriter, r *http.Request, targetURL *url.URL, logger *slog.Logger) {
	// Get the underlying connection
	hijacker, ok := w.(http.Hijacker)
	if !ok {
		logger.Error("response writer doesn't support hijacking")
		http.Error(w, "WebSocket not supported", http.StatusInternalServerError)
		return
	}

	// Build WebSocket URL (ws:// instead of http://)
	wsURL := *targetURL
	if wsURL.Scheme == "https" {
		wsURL.Scheme = "wss"
	} else {
		wsURL.Scheme = "ws"
	}
	wsURL.Path = r.URL.Path
	wsURL.RawQuery = r.URL.RawQuery

	// Connect to Vite's WebSocket
	targetConn, err := net.Dial("tcp", targetURL.Host)
	if err != nil {
		logger.Error("failed to connect to vite websocket", "error", err)
		http.Error(w, "Failed to connect to Vite", http.StatusBadGateway)
		return
	}
	defer targetConn.Close()

	// Hijack the client connection
	clientConn, buf, err := hijacker.Hijack()
	if err != nil {
		logger.Error("failed to hijack connection", "error", err)
		targetConn.Close()
		return
	}
	defer clientConn.Close()

	// Forward the upgrade request to Vite
	err = r.Write(targetConn)
	if err != nil {
		logger.Error("failed to write upgrade request", "error", err)
		return
	}

	// Copy data bidirectionally
	errChan := make(chan error, 2)

	// Copy from client to Vite
	go func() {
		_, err := io.Copy(targetConn, buf)
		errChan <- err
	}()

	// Copy from Vite to client
	go func() {
		_, err := io.Copy(clientConn, targetConn)
		errChan <- err
	}()

	// Wait for either direction to complete
	err = <-errChan
	if err != nil && err != io.EOF {
		logger.Debug("websocket proxy connection closed", "error", err)
	}
}
