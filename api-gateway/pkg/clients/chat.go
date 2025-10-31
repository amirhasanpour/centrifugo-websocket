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

type ChatClient struct {
	conn    *grpc.ClientConn
	service proto.ChatServiceClient
}

func NewChatClient() (*ChatClient, error) {
	chatServiceURL := os.Getenv("CHAT_SERVICE_URL")
	if chatServiceURL == "" {
		chatServiceURL = "localhost:50052"
	}

	conn, err := grpc.Dial(chatServiceURL,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithTimeout(10*time.Second),
	)
	if err != nil {
		return nil, err
	}

	service := proto.NewChatServiceClient(conn)
	log.Printf("Connected to Chat Service at %s", chatServiceURL)

	return &ChatClient{
		conn:    conn,
		service: service,
	}, nil
}

func (c *ChatClient) Close() error {
	return c.conn.Close()
}

func (c *ChatClient) CreateRoom(ctx context.Context, name, description, userID string) (*proto.CreateRoomResponse, error) {
	req := &proto.CreateRoomRequest{
		Name:        name,
		Description: description,
		UserId:      userID,
	}
	return c.service.CreateRoom(ctx, req)
}

func (c *ChatClient) GetRooms(ctx context.Context) (*proto.GetRoomsResponse, error) {
	req := &proto.GetRoomsRequest{}
	return c.service.GetRooms(ctx, req)
}

func (c *ChatClient) JoinRoom(ctx context.Context, roomID, userID string) (*proto.JoinRoomResponse, error) {
	req := &proto.JoinRoomRequest{
		RoomId: roomID,
		UserId: userID,
	}
	return c.service.JoinRoom(ctx, req)
}

func (c *ChatClient) SendMessage(ctx context.Context, roomID, userID, content string) (*proto.SendMessageResponse, error) {
	req := &proto.SendMessageRequest{
		RoomId:  roomID,
		UserId:  userID,
		Content: content,
	}
	return c.service.SendMessage(ctx, req)
}

func (c *ChatClient) GetRoomMessages(ctx context.Context, roomID string, limit int32) (*proto.GetRoomMessagesResponse, error) {
	req := &proto.GetRoomMessagesRequest{
		RoomId: roomID,
		Limit:  limit,
	}
	return c.service.GetRoomMessages(ctx, req)
}