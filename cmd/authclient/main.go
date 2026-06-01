// Command authclient is a tiny demo gRPC client for the auth microservice.
// It registers (or logs in) a user and then verifies the returned token,
// proving the gRPC contract works end-to-end. Run cmd/auth first.
//
// Usage: go run ./cmd/authclient <email> <password>
package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"story-go-mysql/internal/authpb"
)

func main() {
	if err := run(); err != nil {
		slog.Error("authclient failed", "error", err)
		os.Exit(1)
	}
}

func run() error {
	email, password := "grpc-demo@test.com", "password123"
	if len(os.Args) >= 3 {
		email, password = os.Args[1], os.Args[2]
	}

	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()
	client := authpb.NewAuthServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Intenta registrar; si el email ya existe, hace login.
	res, err := client.Register(ctx, &authpb.RegisterRequest{Email: email, Password: password})
	if err != nil {
		fmt.Println("Register falló (¿ya existe?), probando Login…", err)
		res, err = client.Login(ctx, &authpb.LoginRequest{Email: email, Password: password})
		if err != nil {
			return err
		}
	}
	fmt.Printf("OK -> userId=%d email=%s token=%.20s...\n", res.GetUserId(), res.GetEmail(), res.GetAccessToken())

	// Verifica el token recibido.
	v, err := client.Verify(ctx, &authpb.VerifyRequest{AccessToken: res.GetAccessToken()})
	if err != nil {
		return err
	}
	fmt.Printf("Verify -> valid=%v userId=%d\n", v.GetValid(), v.GetUserId())

	// Verifica un token basura (debe dar valid=false).
	bad, err := client.Verify(ctx, &authpb.VerifyRequest{AccessToken: "token-basura"})
	if err != nil {
		return err
	}
	fmt.Printf("Verify (token basura) -> valid=%v\n", bad.GetValid())
	return nil
}
