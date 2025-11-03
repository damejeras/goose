package frontend

import (
	"io/fs"
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

// NewSPAHandlerFS creates an HTTP handler that serves static files from an embedded filesystem.
// If a requested file is not found, it serves index.html to support client-side routing.
func NewSPAHandlerFS(fsys fs.FS) http.Handler {
	// Wrap the filesystem with http.FS for proper serving
	httpFS := http.FS(fsys)
	fileServer := http.FileServer(httpFS)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Clean the path to prevent directory traversal
		path := r.URL.Path
		if path == "/" {
			path = "/index.html"
		}

		// Try to open the file (remove leading slash for fs.FS)
		_, err := fsys.Open(path[1:])
		if err != nil {
			// File doesn't exist, serve index.html
			r.URL.Path = "/index.html"
			fileServer.ServeHTTP(w, r)
			return
		}

		// File exists, serve it normally
		fileServer.ServeHTTP(w, r)
	})
}
