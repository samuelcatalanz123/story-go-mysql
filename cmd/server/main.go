// Command server is the entry point of the Story API. It loads
// configuration, opens the database, wires every layer together (the
// composition root) and runs an HTTP server with graceful shutdown.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/config"
	"story-go-mysql/internal/handler"
	"story-go-mysql/internal/repository"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/storage"
)

func main() {
	if err := run(); err != nil {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	// Refuse to start in production with the publicly known default secret:
	// it would let anyone forge valid JWTs. Railway and similar hosts set PORT.
	if cfg.UsesInsecureJWTSecret() {
		if _, deployed := os.LookupEnv("PORT"); deployed {
			return errors.New("JWT_SECRET must be set in production (refusing to start with the insecure default)")
		}
		slog.Warn("using the insecure default JWT_SECRET; set JWT_SECRET outside local development")
	}

	db, err := storage.NewMySQL(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	// Crea el esquema si no existe (idempotente). Así una base de datos
	// nueva (p. ej. en Railway) queda lista al arrancar.
	migCtx, migCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer migCancel()
	if err := storage.Migrate(migCtx, db); err != nil {
		return err
	}

	// Repositories (data access).
	characterRepo := repository.NewCharacterRepository(db)
	locationRepo := repository.NewLocationRepository(db)
	sceneRepo := repository.NewSceneRepository(db)
	storyRepo := repository.NewStoryRepository(db)
	organizationRepo := repository.NewOrganizationRepository(db)
	conflictRepo := repository.NewConflictRepository(db)

	// Services (business logic).
	characterSvc := service.NewCharacterService(characterRepo, organizationRepo)
	locationSvc := service.NewLocationService(locationRepo)
	sceneSvc := service.NewSceneService(sceneRepo, characterRepo, locationRepo)
	storySvc := service.NewStoryService(storyRepo)
	organizationSvc := service.NewOrganizationService(organizationRepo)
	conflictSvc := service.NewConflictService(conflictRepo)

	// Auth: token manager + user repository + service.
	tokenManager := auth.NewTokenManager(cfg.JWTSecret, 24*time.Hour)
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, tokenManager)

	// Handlers (HTTP) and router.
	router := handler.Router(
		tokenManager,
		handler.NewAuthHandler(authSvc),
		handler.NewCharacterHandler(characterSvc),
		handler.NewLocationHandler(locationSvc),
		handler.NewSceneHandler(sceneSvc),
		handler.NewStoryHandler(storySvc),
		handler.NewOrganizationHandler(organizationSvc),
		handler.NewConflictHandler(conflictSvc),
	)

	// El binario sirve la API en /api/* y el frontend compilado en el resto.
	app := handler.WithFrontend(router, cfg.WebDir)

	server := &http.Server{
		Addr:              cfg.ServerAddr,
		Handler:           app,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// Run the server until an OS signal asks us to stop.
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	serverErr := make(chan error, 1)
	go func() {
		slog.Info("server running", "addr", cfg.ServerAddr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
		}
	}()

	select {
	case err := <-serverErr:
		return err
	case <-ctx.Done():
		slog.Info("shutting down")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return server.Shutdown(shutdownCtx)
	}
}
