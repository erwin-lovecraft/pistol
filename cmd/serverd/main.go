package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/erwin-lovecraft/pistol/internal/adapters/handler"
	"github.com/erwin-lovecraft/pistol/internal/adapters/repository"
	"github.com/erwin-lovecraft/pistol/internal/config"
	"github.com/erwin-lovecraft/pistol/internal/core/services"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("serve exit abnormally: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg := config.ReadFromENV()

	roomRepo := repository.NewInMemoryRoomRepository()
	eventRepo := repository.NewInMemoryEventRepository()
	service := services.NewService(roomRepo, eventRepo)
	hdl := handler.New(service)

	log.Printf("listening on port %s", cfg.Port)
	srv := http.Server{
		Addr:        fmt.Sprintf(":%s", cfg.Port),
		Handler:     routes(hdl),
		ReadTimeout: 5 * time.Second,
		//WriteTimeout: 10 * time.Second, // SSE Endpoint need keep-alive
		IdleTimeout: 2 * time.Minute,
	}
	return srv.ListenAndServe()
}

func routes(hdl handler.Handler) http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(httprate.Limit(
		100,
		1*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint)),
	)

	r.Get("/healthz", healthz)
	r.Get("/rooms/{roomID}/views", hdl.ViewRoom())
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Get("/rooms/{roomID}/events", hdl.ListenEvents())
		v1.Get("/rooms/{roomID}", hdl.ListEvents())
		v1.Handle("/rooms/{roomID}/push", hdl.PushEvent())
	})

	return r
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
