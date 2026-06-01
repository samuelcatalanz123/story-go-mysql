// Command auth is the authentication microservice. It serves the AuthService
// gRPC contract (proto/auth/v1) on :9000, reusing the same business logic as
// the monolith. This demonstrates extracting one domain into its own process.
package main

import (
	"context"
	"log/slog"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"

	"story-go-mysql/internal/auth"
	"story-go-mysql/internal/authpb"
	"story-go-mysql/internal/config"
	"story-go-mysql/internal/model"
	"story-go-mysql/internal/repository"
	"story-go-mysql/internal/service"
	"story-go-mysql/internal/storage"
)

func main() {
	if err := run(); err != nil {
		slog.Error("auth service failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := config.Load()

	db, err := storage.NewMySQL(cfg.DB)
	if err != nil {
		return err
	}
	defer db.Close()

	migCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := storage.Migrate(migCtx, db); err != nil {
		return err
	}

	tokenManager := auth.NewTokenManager(cfg.JWTSecret, 15*time.Minute)
	authSvc := service.NewAuthService(
		repository.NewUserRepository(db),
		repository.NewRefreshTokenRepository(db),
		tokenManager,
		30*24*time.Hour,
	)

	addr := ":9000"
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, &authServer{svc: authSvc, tokens: tokenManager})

	slog.Info("auth gRPC service running", "addr", addr)
	return grpcServer.Serve(lis)
}

// authServer implements the generated AuthServiceServer interface by delegating
// to the existing AuthService and TokenManager.
type authServer struct {
	authpb.UnimplementedAuthServiceServer
	svc    *service.AuthService
	tokens *auth.TokenManager
}

func (s *authServer) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.AuthResponse, error) {
	res, _, err := s.svc.Register(ctx, model.RegisterRequest{Email: req.GetEmail(), Password: req.GetPassword()})
	if err != nil {
		return nil, err
	}
	return &authpb.AuthResponse{AccessToken: res.Token, UserId: res.User.ID, Email: res.User.Email}, nil
}

func (s *authServer) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.AuthResponse, error) {
	res, _, err := s.svc.Login(ctx, model.LoginRequest{Email: req.GetEmail(), Password: req.GetPassword()})
	if err != nil {
		return nil, err
	}
	return &authpb.AuthResponse{AccessToken: res.Token, UserId: res.User.ID, Email: res.User.Email}, nil
}

func (s *authServer) Verify(_ context.Context, req *authpb.VerifyRequest) (*authpb.VerifyResponse, error) {
	userID, err := s.tokens.Parse(req.GetAccessToken())
	if err != nil {
		// Token inválido/expirado: no es un error gRPC, simplemente valid=false.
		return &authpb.VerifyResponse{Valid: false}, nil
	}
	return &authpb.VerifyResponse{Valid: true, UserId: userID}, nil
}
