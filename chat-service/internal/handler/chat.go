package handler

import (
	"context"
	"time"

	"chat-service/internal/proto"
	"chat-service/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatHandler struct {
	proto.UnimplementedChatServiceServer
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

func (h *ChatHandler) CreateRoom(ctx context.Context, req *proto.CreateRoomRequest) (*proto.CreateRoomResponse, error) {
	room, err := h.chatService.CreateRoom(req.Name, req.Description, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.CreateRoomResponse{
		Room: &proto.Room{
			Id:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			CreatedBy:   room.CreatedBy,
			CreatedAt:   room.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   room.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *ChatHandler) GetRooms(ctx context.Context, req *proto.GetRoomsRequest) (*proto.GetRoomsResponse, error) {
	rooms, err := h.chatService.GetRooms()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoRooms := make([]*proto.Room, len(rooms))
	for i, room := range rooms {
		protoRooms[i] = &proto.Room{
			Id:          room.ID,
			Name:        room.Name,
			Description: room.Description,
			CreatedBy:   room.CreatedBy,
			CreatedAt:   room.CreatedAt.Format(time.RFC3339),
			UpdatedAt:   room.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetRoomsResponse{
		Rooms: protoRooms,
	}, nil
}

func (h *ChatHandler) JoinRoom(ctx context.Context, req *proto.JoinRoomRequest) (*proto.JoinRoomResponse, error) {
	err := h.chatService.JoinRoom(req.RoomId, req.UserId)
	if err != nil {
		if err == service.ErrRoomNotFound {
			return nil, status.Error(codes.NotFound, "room not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.JoinRoomResponse{
		Success: true,
		RoomId:  req.RoomId,
	}, nil
}

func (h *ChatHandler) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	// Note: We need username here, but it's not in the request
	// For now, we'll use user_id as username, but in real implementation
	// we would need to get username from auth service
	username := req.UserId // Temporary - should be fetched from user service

	message, err := h.chatService.SendMessage(req.RoomId, req.UserId, username, req.Content)
	if err != nil {
		if err == service.ErrRoomNotFound {
			return nil, status.Error(codes.NotFound, "room not found")
		}
		if err == service.ErrNotRoomMember {
			return nil, status.Error(codes.PermissionDenied, "user is not a member of this room")
		}
		if err == service.ErrInvalidMessage {
			return nil, status.Error(codes.InvalidArgument, "invalid message content")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &proto.SendMessageResponse{
		Message: &proto.Message{
			Id:        message.ID,
			RoomId:    message.RoomID,
			UserId:    message.UserID,
			Username:  message.Username,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (h *ChatHandler) GetRoomMessages(ctx context.Context, req *proto.GetRoomMessagesRequest) (*proto.GetRoomMessagesResponse, error) {
	messages, err := h.chatService.GetRoomMessages(req.RoomId, int(req.Limit))
	if err != nil {
		if err == service.ErrRoomNotFound {
			return nil, status.Error(codes.NotFound, "room not found")
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoMessages := make([]*proto.Message, len(messages))
	for i, message := range messages {
		protoMessages[i] = &proto.Message{
			Id:        message.ID,
			RoomId:    message.RoomID,
			UserId:    message.UserID,
			Username:  message.Username,
			Content:   message.Content,
			CreatedAt: message.CreatedAt.Format(time.RFC3339),
		}
	}

	return &proto.GetRoomMessagesResponse{
		Messages: protoMessages,
	}, nil
}