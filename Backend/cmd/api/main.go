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
	"github.com/WBetEngine/Seosementara/Backend/internal/crypto"
	"github.com/WBetEngine/Seosementara/Backend/internal/handler"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/facebook"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/service"
	"github.com/WBetEngine/Seosementara/Backend/internal/pixel/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := config.Load()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	st, encKey := openStore(cfg, logger)
	capi := facebook.NewCAPIClient()
	fbSvc := service.NewFacebookService(st, capi, encKey)

	tpl, err := handler.LoadTemplates(cfg.AdminTemplatesDir)
	if err != nil {
		logger.Warn("templates", "err", err)
	}

	adminFB := &handler.AdminPixelFacebook{Svc: fbSvc, Templates: tpl}
	collect := &handler.CollectHandler{FB: fbSvc}

	r := chi.NewRouter()
	r.Use(middleware.RequestID, middleware.RealIP, middleware.Logger, middleware.Recoverer)

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

func openStore(cfg config.Config, logger *slog.Logger) (store.PixelStore, []byte) {
	if cfg.DatabaseURL != "" {
		pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
		if err != nil {
			logger.Error("postgres", "err", err)
			os.Exit(1)
		}
		logger.Info("store", "type", "postgres")
		key, err := resolveEncKey(cfg.PixelEncryptionKey)
		if err != nil {
			logger.Warn("encryption", "err", err, "msg", "dev key generated")
			key, _ = devKey()
		}
		return store.NewPostgresStore(pool), key
	}
	logger.Info("store", "type", "memory", "hint", "set DATABASE_URL for production")
	key, _ := devKey()
	if cfg.PixelEncryptionKey != "" {
		if k, err := crypto.KeyFromEnv(cfg.PixelEncryptionKey); err == nil {
			key = k
		}
	}
	return store.NewMemoryStore(), key
}

func resolveEncKey(b64 string) ([]byte, error) {
	if b64 != "" {
		return crypto.KeyFromEnv(b64)
	}
	return devKey()
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
