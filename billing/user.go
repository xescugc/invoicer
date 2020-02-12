package billing

import (
	"context"
	"errors"
	"fmt"

	"github.com/xescugc/invoicer/user"
)

func (b *billing) CreateUser(ctx context.Context, u *user.User) error {
	nu, err := b.users.Find(ctx)
	if err != nil && !errors.Is(err, ErrNotFoundUser) {
		return fmt.Errorf("could not find User: %w", err)
	}

	if nu != nil {
		return ErrAlreadyExistsUser
	}

	err = b.users.Create(ctx, u)
	if err != nil {
		return fmt.Errorf("could not create User: %w", err)
	}

	return nil
}

func (b *billing) GetUser(ctx context.Context) (*user.User, error) {
	u, err := b.users.Find(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get User: %w", err)
	}

	return u, nil
}

func (b *billing) UpdateUser(ctx context.Context, u *user.User) error {
	err := b.users.Update(ctx, u)
	if err != nil {
		return fmt.Errorf("could not update User: %w", err)
	}

	return nil
}
