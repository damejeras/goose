package frontend

import (
	"net/http"
	"os"
	"path/filepath"
)

// Handler returns an HTTP handler that serves the SPA from the dist directory.
// It serves index.html for any routes that don't match existing files.
func Handler() http.Handler {
	return NewSPAHandler("frontend/dist")
}

// NewSPAHandler creates an HTTP handler that serves static files from the dist directory.
// If a requested file is not found, it serves index.html to support client-side routing.
func NewSPAHandler(distPath string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Build the full file path
		path := filepath.Join(distPath, r.URL.Path)

		// Check if the file exists
		info, err := os.Stat(path)
		if os.IsNotExist(err) || info.IsDir() {
			// File doesn't exist or is a directory, serve index.html
			http.ServeFile(w, r, filepath.Join(distPath, "index.html"))
			return
		}

		if err != nil {
			// Other error occurred
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// File exists, serve it
		http.ServeFile(w, r, path)
	})
}
