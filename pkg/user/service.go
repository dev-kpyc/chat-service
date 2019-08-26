package user

import "context"

type User struct {
	ID   int64
	Name string
}

type Repository interface {
	StoreUser(ctx context.Context, name string) (userID int64, err error)
	GetUser(ctx context.Context, userID int64) (*User, error)
}

type Service struct {
	repo Repository
}

// New ...
func New(repo Repository) Service {
	return Service{repo}
}

func (svc *Service) CreateUser(ctx context.Context, name string) (userID int64, err error) {

	userID, err = svc.repo.StoreUser(ctx, name)

	if err != nil {
		return 0, err
	}

	return
}

func (svc *Service) GetUser(ctx context.Context, userID int64) (*User, error) {

	return svc.repo.GetUser(ctx, userID)
}
