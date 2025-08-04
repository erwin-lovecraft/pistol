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

func (h Handler) CreateRoom() http.HandlerFunc {
	type request struct {
		Name   string `json:"name"`
		Avatar string `json:"avatar"`
	}

	type response struct {
		ID   string `json:"id"`
		Link string `json:"link"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var req request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		room, link, err := h.svc.CreateRoom(r.Context(), req.Name, req.Avatar)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response{
			ID:   room.ID,
			Link: link,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (h Handler) ListRoom() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rooms, err := h.svc.ListRoom(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(rooms); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
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
