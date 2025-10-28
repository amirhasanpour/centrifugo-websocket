package clients

import (
	"context"
	"log"
	"os"
	"time"

	"api-gateway/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthClient struct {
	conn    *grpc.ClientConn
	service proto.AuthServiceClient
}

func NewAuthClient() (*AuthClient, error) {
	authServiceURL := os.Getenv("AUTH_SERVICE_URL")
	if authServiceURL == "" {
		authServiceURL = "localhost:50051"
	}

	conn, err := grpc.Dial(authServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	service := proto.NewAuthServiceClient(conn)
	log.Printf("Connected to Auth Service at %s", authServiceURL)

	return &AuthClient{
		conn:    conn,
		service: service,
	}, nil
}

func (c *AuthClient) Close() error {
	return c.conn.Close()
}

func (c *AuthClient) Register(ctx context.Context, username, email, password string) (*proto.RegisterResponse, error) {
	req := &proto.RegisterRequest{
		Username: username,
		Email:    email,
		Password: password,
	}
	return c.service.Register(ctx, req)
}

func (c *AuthClient) Login(ctx context.Context, email, password string) (*proto.LoginResponse, error) {
	req := &proto.LoginRequest{
		Email:    email,
		Password: password,
	}
	return c.service.Login(ctx, req)
}

func (c *AuthClient) ValidateToken(ctx context.Context, token string) (*proto.ValidateTokenResponse, error) {
	req := &proto.ValidateTokenRequest{
		Token: token,
	}
	return c.service.ValidateToken(ctx, req)
}

func (c *AuthClient) GetUser(ctx context.Context, userID string) (*proto.GetUserResponse, error) {
	req := &proto.GetUserRequest{
		UserId: userID,
	}
	return c.service.GetUser(ctx, req)
}