package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

func loadTemplates(dir string) (*template.Template, error) {
	pattern := filepath.Join(dir, "*.html")
	glob, err := template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("failed to load templates: %w", err)
	}

	return glob, nil
}

func deriveBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	} else if proto := r.Header.Get("X-Forwarded-Proto"); strings.EqualFold(proto, "https") {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}
	// ensure trailing slash
	return scheme + "://" + host + "/"
}
