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
	"github.com/erwin-lovecraft/pistol/internal/web"
	"github.com/erwin-lovecraft/pistol/migrations"
	pkgmiddleware "github.com/erwin-lovecraft/pistol/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("serve exit abnormally: %v", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	cfg := config.ReadFromENV()

	// Setup DB connections
	dbPool, err := pgxpool.New(ctx, cfg.PGURL)
	if err != nil {
		return fmt.Errorf("create db connection: %w", err)
	}
	defer dbPool.Close()

	// Setup migrations
	goose.SetBaseFS(migrations.MigrationFS)
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("goose set dialect: %w", err)
	}
	if err := goose.UpContext(ctx, stdlib.OpenDBFromPool(dbPool), "."); err != nil {
		return fmt.Errorf("migrate up: %w", err)
	}

	// Setup ID generator
	if err := repository.SetupIDGenerator(); err != nil {
		return err
	}

	// DI settings
	roomRepo := repository.NewInMemoryRoomRepository()
	eventRepo := repository.NewEventRepository(dbPool)
	service := services.NewService(roomRepo, eventRepo)
	hdl, err := handler.New(service)
	if err != nil {
		return err
	}

	// Start server
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
	r.Use(httprate.Limit(
		100,
		1*time.Second,
		httprate.WithKeyFuncs(httprate.KeyByIP, httprate.KeyByEndpoint)),
	)

	r.Get("/healthz", healthz)
	r.Get("/", hdl.Home())
	r.Handle("/*", http.FileServer(http.FS(web.FS)))
	r.Get("/rooms/{roomID}/views", hdl.ViewRoom())
	r.Route("/api/v1", func(v1 chi.Router) {
		v1.Get("/rooms/{roomID}/events", hdl.ListenEvents())
		v1.Get("/rooms/{roomID}", hdl.ListEvents())
		v1.Handle("/rooms/{roomID}/push", pkgmiddleware.AuthKey(hdl.PushEvent()))
	})

	return r
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
