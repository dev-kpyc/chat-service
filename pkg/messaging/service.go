package messaging

import (
	"context"
	"fmt"

	"github.com/dev-kpyc/chat-service/pkg/user"
)

// Service performs messaging use cases
type Service struct {
	repo  Repository
	users user.Service
}

// Repository defines the repository used to store and retrieve messages
type Repository interface {
	GetUserIDsInChatRoom(ctx context.Context, roomID int64) ([]int64, error)
	SaveChatMessage(ctx context.Context, senderID int64, roomID int64, message string) (int64, error)
	GetChatMessages(ctx context.Context, roomID int64) ([]*ChatMessage, error)
}

// User defines the user making the request
type User struct {
	ID   int64
	Name string
}

// ChatMessage defines a message sent to a chat room
type ChatMessage struct {
	SenderID int64
	Message  string
}

func New(repo Repository, users user.Service) Service {
	return Service{repo, users}
}

func (svc *Service) SendMessage(ctx context.Context, userID int64, roomID int64, message string) (int64, error) {

	_, err := svc.users.GetUser(ctx, userID)

	if err != nil {
		return 0, err
	}

	if err := svc.checkUserInChatRoom(ctx, userID, roomID); err != nil {
		return 0, err
	}

	id, err := svc.repo.SaveChatMessage(ctx, userID, roomID, message)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (svc *Service) GetMessages(ctx context.Context, userID int64, roomID int64) ([]*ChatMessage, error) {

	_, err := svc.users.GetUser(ctx, userID)

	if err != nil {
		return nil, err
	}

	if err := svc.checkUserInChatRoom(ctx, userID, roomID); err != nil {
		return nil, err
	}

	messages, err := svc.repo.GetChatMessages(ctx, roomID)

	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (svc *Service) checkUserInChatRoom(ctx context.Context, userID int64, roomID int64) error {

	userIDs, err := svc.repo.GetUserIDsInChatRoom(ctx, roomID)

	if err != nil {
		return err
	}

	if !contains(userID, userIDs) {
		return fmt.Errorf("User does not belong to this chat room")
	}

	return nil
}

func contains(userID int64, userIDs []int64) bool {
	for _, u := range userIDs {
		if u == userID {
			return true
		}
	}
	return false
}
