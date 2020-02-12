package customer

import "context"

//go:generate mockgen -destination=../mock/customer_repository.go -mock_names=Repository=CustomerRepository -package mock github.com/xescugc/invoicer/customer Repository

type Repository interface {
	Create(ctx context.Context, c *Customer) error
	Find(ctx context.Context, canonical string) (*Customer, error)
	Filter(ctx context.Context) ([]*Customer, error)
	Update(ctx context.Context, ID string, c *Customer) error
	Delete(ctx context.Context, ID string) error
}
