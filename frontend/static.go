package frontend

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var staticFiles embed.FS

func Handler() http.Handler {
	staticFS, err := fs.Sub(staticFiles, "dist")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(staticFS))
}
