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
	"story-go-mysql/internal/cache"
	"story-go-mysql/internal/config"
	"story-go-mysql/internal/email"
	"story-go-mysql/internal/handler"
	"story-go-mysql/internal/oauth"
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
	refreshRepo := repository.NewRefreshTokenRepository(db)
	oauthRepo := repository.NewOAuthAccountRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// Cache: Redis when reachable, otherwise a no-op so the app keeps working.
	var appCache cache.Cache = cache.Noop{}
	if cfg.Redis.Addr != "" {
		if rc, cacheErr := cache.NewRedis(context.Background(), cfg.Redis.Addr, cfg.Redis.Password); cacheErr != nil {
			slog.Warn("redis unavailable; caching and rate limiting disabled", "error", cacheErr)
		} else {
			appCache = rc
			slog.Info("redis enabled", "addr", cfg.Redis.Addr)
		}
	}

	// Services (business logic).
	characterSvc := service.NewCharacterService(characterRepo, organizationRepo, appCache)
	locationSvc := service.NewLocationService(locationRepo)
	sceneSvc := service.NewSceneService(sceneRepo, characterRepo, locationRepo)
	storySvc := service.NewStoryService(storyRepo)
	organizationSvc := service.NewOrganizationService(organizationRepo)
	conflictSvc := service.NewConflictService(conflictRepo)

	// Auth: short-lived access token (15 min) + long-lived refresh token (30
	// days) stored as a hash in the DB and sent to the client in an HttpOnly
	// cookie.
	const refreshTTL = 30 * 24 * time.Hour
	tokenManager := auth.NewTokenManager(cfg.JWTSecret, 15*time.Minute)
	userRepo := repository.NewUserRepository(db)
	authSvc := service.NewAuthService(userRepo, refreshRepo, tokenManager, refreshTTL)

	// "Sign in with Google" is optional: only enabled when credentials are set.
	var googleAuth *oauth.GoogleAuthenticator
	if cfg.Google.ClientID != "" {
		googleAuth, err = oauth.NewGoogleAuthenticator(
			context.Background(), cfg.Google.ClientID, cfg.Google.ClientSecret, cfg.Google.RedirectURI)
		if err != nil {
			return err
		}
		slog.Info("google login enabled")
	} else {
		slog.Info("google login disabled (set GOOGLE_CLIENT_ID to enable)")
	}
	// A nil *GoogleAuthenticator must be passed as a nil interface, so branch.
	var oauthSvc *service.OAuthService
	if googleAuth != nil {
		oauthSvc = service.NewOAuthService(userRepo, oauthRepo, googleAuth, tokenManager, refreshRepo, refreshTTL)
	} else {
		oauthSvc = service.NewOAuthService(userRepo, oauthRepo, nil, tokenManager, refreshRepo, refreshTTL)
	}

	// Email sender: log-only in development, real SMTP when SMTP_HOST is set.
	var mailer email.Sender = email.LogSender{}
	if cfg.SMTP.Host != "" {
		mailer = email.SMTPSender{Addr: cfg.SMTP.Host + ":" + cfg.SMTP.Port, From: cfg.SMTP.From}
		slog.Info("email via SMTP", "addr", cfg.SMTP.Host+":"+cfg.SMTP.Port)
	} else {
		slog.Info("email disabled (logging only; set SMTP_HOST to send)")
	}
	passwordResetSvc := service.NewPasswordResetService(userRepo, passwordResetRepo, mailer, cfg.AppBaseURL, time.Hour)

	// Handlers (HTTP) and router.
	router := handler.Router(
		tokenManager,
		handler.NewAuthHandler(authSvc, oauthSvc, refreshTTL),
		handler.NewPasswordHandler(passwordResetSvc),
		handler.NewCharacterHandler(characterSvc, cfg.UploadDir),
		handler.NewLocationHandler(locationSvc, cfg.UploadDir),
		handler.NewSceneHandler(sceneSvc),
		handler.NewStoryHandler(storySvc),
		handler.NewOrganizationHandler(organizationSvc),
		handler.NewConflictHandler(conflictSvc),
		cfg.UploadDir,
		appCache,
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
