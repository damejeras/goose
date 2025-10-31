package main

import (
	"log"
	"net/http"

	"github.com/damejeras/goose/api/v1/v1connect"
	"github.com/damejeras/goose/frontend"
	"github.com/damejeras/goose/internal/greeter"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	mux := http.NewServeMux()

	path, handler := v1connect.NewGreeterServiceHandler(&greeter.Server{})
	mux.Handle(path, handler)

	// Serve static files from frontend package
	mux.Handle("/", frontend.Handler())

	log.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{})); err != nil {
		log.Fatal(err)
	}
}
