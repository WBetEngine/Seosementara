package main

import (
	"context"
	"crypto/rand"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/WBetEngine/Seosementara/Backend/internal/config"
	"github.com/WBetEngine/Seosementara/Backend/internal/handler"
	authmw "github.com/WBetEngine/Seosementara/Backend/internal/middleware"
	"github.com/WBetEngine/Seosementara/Backend/internal/migrate"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/facebook"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/service"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/store"
	setupservice "github.com/WBetEngine/Seosementara/Backend/internal/setup/service"
	setupstore "github.com/WBetEngine/Seosementara/Backend/internal/setup/store"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	encKey, err := cfg.ResolveEncryptionKey()
	if err != nil {
		logger.Warn("encryption", "err", err, "msg", "dev key")
		encKey, _ = devKey()
	}

	var pool *pgxpool.Pool
	if cfg.DatabaseURL != "" {
		pool, err = pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			logger.Error("postgres", "err", err)
			os.Exit(1)
		}
		logger.Info("database", "type", "postgres")
		migDir := os.Getenv("MIGRATIONS_DIR")
		if migDir == "" {
			migDir = "migrations"
		}
		if err := migrate.Up(context.Background(), pool, migDir); err != nil {
			logger.Error("migrate", "err", err)
			os.Exit(1)
		}
		logger.Info("migrate", "status", "ok", "dir", migDir)
	}

	var pixelSt store.PixelStore = store.NewMemoryStore()
	var setupSt setupstore.SetupStore = setupstore.NewMemoryStore()
	if pool != nil {
		pixelSt = store.NewPostgresStore(pool)
		setupSt = setupstore.NewPostgresStore(pool)
	} else {
		logger.Info("database", "type", "memory", "hint", "set DATABASE_URL for production")
	}

	capi := facebook.NewCAPIClient()
	fbSvc := service.NewFacebookService(pixelSt, capi, encKey)
	setupSvc := setupservice.NewSetupService(setupSt, encKey)

	tpl, err := handler.LoadTemplates(cfg.AdminTemplatesDir)
	if err != nil {
		logger.Warn("templates", "err", err)
	}

	adminFB := &handler.AdminPixelFacebook{Svc: fbSvc, Templates: tpl}
	adminCF := &handler.AdminCloudflare{Svc: setupSvc, PartialsDir: cfg.AdminPartialsDir}
	collect := &handler.CollectHandler{FB: fbSvc}

	r := chi.NewRouter()
	r.Use(chimw.RequestID, chimw.RealIP, chimw.Logger, chimw.Recoverer)

	r.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Post("/collect", collect.ServeHTTP)
	r.Get("/sseo-track.js", handler.ServeTrackScript(cfg.StaticDir))
	fileServer := http.FileServer(http.Dir(cfg.StaticDir))
	r.Handle("/admin/static/*", http.StripPrefix("/admin/static/", fileServer))

	r.Route("/admin/pixel/facebook", adminFB.Routes)

	r.Route("/api/admin/pixel/facebook", func(r chi.Router) {
		r.Get("/setup", adminFB.APIGetSetup)
		r.Post("/setup", adminFB.APIPostSetup)
		r.Post("/test-connection", adminFB.APITestConnection)
		r.Post("/test-event", adminFB.APITestEvent)
		r.Get("/diagnostics", adminFB.APIDiagnostics)
		r.Get("/events", adminFB.APIEvents)
		r.Post("/domains/assign", adminFB.APIAssignDomain)
	})

	guard := authmw.RequireSuperAdmin
	handler.MountCloudflare(r, adminCF, guard)

	// background dispatcher
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go runDispatcher(ctx, fbSvc, cfg.DispatchIntervalSec, logger)

	srv := &http.Server{Addr: cfg.Addr, Handler: r}
	go func() {
		logger.Info("api listening", "addr", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}

func devKey() ([]byte, error) {
	k := make([]byte, 32)
	_, err := rand.Read(k)
	return k, err
}

func runDispatcher(ctx context.Context, fb *service.FacebookService, intervalSec int, logger *slog.Logger) {
	t := time.NewTicker(time.Duration(intervalSec) * time.Second)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			sent, failed, err := fb.DispatchPending(ctx, 50)
			if err != nil {
				logger.Error("pixel_dispatch", "err", err)
				continue
			}
			if sent > 0 || failed > 0 {
				logger.Info("pixel_dispatch", "sent", sent, "failed", failed)
			}
		}
	}
}
