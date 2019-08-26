package chatroom

import (
	"context"
	"fmt"

	"github.com/dev-kpyc/chat-service/pkg/user"
)

// ChatRoom ...
type ChatRoom struct {
	ID      int64
	Name    string
	OwnerID int64
}

// Repository ...
type Repository interface {
	GetChatRoom(ctx context.Context, roomID int64) (*ChatRoom, error)
	StoreChatRoom(ctx context.Context, name string, ownerID int64) (int64, error)
	AddUserToChatRoom(ctx context.Context, userID int64, roomID int64) error
	RemoveUserFromChatRoom(ctx context.Context, userID int64, roomID int64) error
	ListChatRooms(ctx context.Context, userID int64) ([]*ChatRoom, error)
	ListOwnedChatRooms(ctx context.Context, ownerID int64) ([]*ChatRoom, error)
}

// New ...
func New(repo Repository, users user.Service) Service {
	return Service{repo, users}
}

// Service ...
type Service struct {
	repo  Repository
	users user.Service
}

// User ...
type User struct {
	ID   int64
	Name string
}

// Create ...
func (svc *Service) Create(ctx context.Context, userID int64, name string) (int64, error) {

	_, err := svc.users.GetUser(ctx, userID)

	if err != nil {
		return 0, err
	}

	id, err := svc.repo.StoreChatRoom(ctx, name, userID)

	if err != nil {
		return 0, err
	}

	err = svc.Join(ctx, userID, id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

// List ...
func (svc *Service) List(ctx context.Context, userID int64) ([]*ChatRoom, error) {

	user, err := svc.users.GetUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	chatrooms, err := svc.repo.ListChatRooms(ctx, user.ID)

	if err != nil {
		return nil, err
	}

	return chatrooms, nil
}

// ListOwned ...
func (svc *Service) ListOwned(ctx context.Context, ownerID int64) ([]*ChatRoom, error) {
	_, err := svc.users.GetUser(ctx, ownerID)

	if err != nil {
		return nil, err
	}

	chatrooms, err := svc.repo.ListOwnedChatRooms(ctx, ownerID)

	if err != nil {
		return nil, err
	}

	return chatrooms, nil
}

// Join ...
func (svc *Service) Join(ctx context.Context, userID int64, roomID int64) error {

	user, err := svc.users.GetUser(ctx, userID)

	if err != nil {
		return err
	}

	return svc.repo.AddUserToChatRoom(ctx, user.ID, roomID)
}

// Leave ...
func (svc *Service) Leave(ctx context.Context, userID int64, roomID int64) error {

	user, err := svc.users.GetUser(ctx, userID)

	if err != nil {
		return err
	}

	chatroom, err := svc.repo.GetChatRoom(ctx, roomID)

	if err != nil {
		return err
	}

	if chatroom.OwnerID == userID {
		return fmt.Errorf("cannot leave chat room because you are the owner")
	}

	return svc.repo.RemoveUserFromChatRoom(ctx, user.ID, roomID)
}
