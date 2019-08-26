package grpc

import (
	"context"
	"log"
	"net"

	pb "github.com/dev-kpyc/chat-service/api"
	"github.com/dev-kpyc/chat-service/pkg/chatroom"
	"github.com/dev-kpyc/chat-service/pkg/messaging"
	"github.com/dev-kpyc/chat-service/pkg/user"
	"google.golang.org/grpc"
)

// Clients is the adapter for grpc endpoints
type Clients struct {
	server    *grpc.Server
	chatroom  *chatroom.Service
	messaging *messaging.Service
	users     *user.Service
}

func New(chatroom *chatroom.Service, messaging *messaging.Service, users *user.Service) *Clients {
	return &Clients{grpc.NewServer(), chatroom, messaging, users}
}

func (c *Clients) RegisterEndpoints(lis net.Listener) {

	pb.RegisterChatRoomServer(c.server, c)

	pb.RegisterUserServer(c.server, c)

	pb.RegisterMessagingServer(c.server, c)

	go c.server.Serve(lis)
}

// CreateUser ...
func (c *Clients) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {

	userID, err := c.users.CreateUser(ctx, req.Name)

	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{UserID: userID}, nil
}

// GetUser ...
func (c *Clients) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {

	user, err := c.users.GetUser(ctx, req.UserID)

	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{Name: user.Name}, nil
}

// CreateChatRoom ...
func (c *Clients) CreateChatRoom(ctx context.Context, req *pb.CreateChatRoomRequest) (*pb.CreateChatRoomResponse, error) {

	roomID, err := c.chatroom.Create(ctx, req.UserID, req.ChatRoomName)

	if err != nil {
		return nil, err
	}

	return &pb.CreateChatRoomResponse{ChatRoomID: roomID}, nil
}

// ListChatRooms ...
func (c *Clients) ListChatRooms(ctx context.Context, req *pb.ListChatRoomsRequest) (*pb.ListChatRoomsResponse, error) {

	chatrooms, err := c.chatroom.List(ctx, req.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	var list []*pb.ListChatRoom
	for _, chatroom := range chatrooms {
		list = append(list, &pb.ListChatRoom{Id: chatroom.ID, Name: chatroom.Name})
	}

	return &pb.ListChatRoomsResponse{Chatrooms: list}, nil
}

// ListOwnedChatRooms ...
func (c *Clients) ListOwnedChatRooms(ctx context.Context, req *pb.ListChatRoomsRequest) (*pb.ListChatRoomsResponse, error) {

	chatrooms, err := c.chatroom.ListOwned(ctx, req.UserID)

	if err != nil {
		return nil, err
	}

	var list []*pb.ListChatRoom
	for _, chatroom := range chatrooms {
		list = append(list, &pb.ListChatRoom{Id: chatroom.ID, Name: chatroom.Name})
	}

	return &pb.ListChatRoomsResponse{Chatrooms: list}, nil
}

// JoinChatRoom ...
func (c *Clients) JoinChatRoom(ctx context.Context, req *pb.JoinChatRoomRequest) (*pb.JoinChatRoomResponse, error) {

	err := c.chatroom.Join(ctx, req.UserID, req.ChatRoomID)

	if err != nil {
		return nil, err
	}

	return &pb.JoinChatRoomResponse{}, nil
}

// LeaveChatRoom ...
func (c *Clients) LeaveChatRoom(ctx context.Context, req *pb.LeaveChatRoomRequest) (*pb.LeaveChatRoomResponse, error) {

	err := c.chatroom.Leave(ctx, req.UserID, req.ChatRoomID)

	if err != nil {
		return nil, err
	}

	return &pb.LeaveChatRoomResponse{}, nil
}

func (c *Clients) SendMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageResponse, error) {

	_, err := c.messaging.SendMessage(ctx, req.UserID, req.ChatRoomID, req.Message)

	if err != nil {
		return nil, err
	}

	return &pb.SendMessageResponse{}, nil
}

func (c *Clients) GetMessages(ctx context.Context, req *pb.GetMessagesRequest) (*pb.GetMessagesResponse, error) {

	messages, err := c.messaging.GetMessages(ctx, req.UserID, req.ChatRoomID)

	if err != nil {
		return nil, err
	}

	log.Println(messages)

	var list []*pb.Message
	for _, message := range messages {
		list = append(list, &pb.Message{SenderID: message.SenderID, Message: message.Message})
	}

	return &pb.GetMessagesResponse{Messages: list}, nil
}
