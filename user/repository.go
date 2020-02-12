package user

import "context"

//go:generate mockgen -destination=../mock/user_repository.go -mock_names=Repository=UserRepository -package mock github.com/xescugc/invoicer/user Repository

type Repository interface {
	Create(ctx context.Context, u *User) error
	Find(ctx context.Context) (*User, error)
	Update(ctx context.Context, u *User) error
}
