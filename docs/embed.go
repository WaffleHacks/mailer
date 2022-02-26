package docs

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

//go:embed gen/*
var embedded embed.FS

// Handler creates a new HTTP handler serving the documentation
func Handler() (http.Handler, error) {
	subFs, err := fs.Sub(embedded, "gen")
	if err != nil {
		return nil, err
	}

	httpFs := http.FileServer(http.FS(subFs))

	// HTTP handler adapted from https://github.com/go-chi/chi/blob/v5.0.7/_examples/fileserver/main.go
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		stripped := http.StripPrefix(pathPrefix, httpFs)
		stripped.ServeHTTP(w, r)
	}), nil
}
