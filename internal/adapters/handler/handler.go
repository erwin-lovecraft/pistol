package handler

import (
	"encoding/json"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/erwin-lovecraft/pistol/internal/core/domain"
	"github.com/erwin-lovecraft/pistol/internal/core/services"
	"github.com/go-chi/chi/v5"
)

var (
	tpl          *template.Template
	loadViewSync sync.Once
)

type Handler struct {
	svc services.Service
}

func New(svc services.Service) Handler {
	return Handler{
		svc: svc,
	}
}

func (h Handler) ListenEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			http.Error(w, "roomID is required", http.StatusBadRequest)
			return
		}

		cl, err := h.svc.ListenEvents(r.Context(), roomID, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		<-cl.Wait()
	}
}

func (h Handler) ListEvents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			http.Error(w, "roomID is required", http.StatusBadRequest)
			return
		}

		var pagination Pagination
		if err := pagination.FromRequest(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		rs, hasMore, err := h.svc.ListEvents(r.Context(), roomID, pagination.Page, pagination.Size)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"data":    rs,
			"page":    pagination.Page,
			"size":    pagination.Size,
			"hasMore": hasMore,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h Handler) PushEvent() http.HandlerFunc {
	type response struct {
		Message string `json:"message"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			http.Error(w, "roomID is required", http.StatusBadRequest)
			return
		}

		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := h.svc.PushEvent(r.Context(), roomID, domain.Event{
			Method:      r.Method,
			Header:      r.Header,
			QueryParams: r.URL.Query(),
			Body:        reqBody,
		}); err != nil {
			http.Error(w, "relay error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response{Message: "ok"}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h Handler) ViewRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		roomID := chi.URLParam(r, "roomID")
		if roomID == "" {
			http.Error(w, "roomID is required", http.StatusBadRequest)
			return
		}

		// Load template once
		loadTemplates("internal/web")

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tpl.ExecuteTemplate(w, "view.html", map[string]string{
			"RoomID": roomID,
		}); err != nil {
			http.Error(w, "failed to render template", http.StatusInternalServerError)
		}
	}
}

func loadTemplates(dir string) {
	loadViewSync.Do(func() {
		pattern := filepath.Join(dir, "*.html")
		globTpl, err := template.ParseGlob(pattern)
		if err != nil {
			panic(err)
		}

		tpl = globTpl
	})
}
