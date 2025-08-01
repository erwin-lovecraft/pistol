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
	service := services.NewService(roomRepo)
	hdl := handler.New(service)

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

	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Post("/rooms", hdl.CreateRoom())
		v1.Get("/rooms/{roomID}/events", hdl.ListenEvents())
		v1.Post("/rooms/{roomID}/relay", hdl.Relay())
	})

	return r
}
