package invoice

import "context"

//go:generate mockgen -destination=../mock/invoice_repository.go -mock_names=Repository=InvoiceRepository -package mock github.com/xescugc/invoicer/invoice Repository

type Repository interface {
	Create(ctx context.Context, i *Invoice) error
	Find(ctx context.Context, number string) (*Invoice, error)
	Filter(ctx context.Context) ([]*Invoice, error)
	Update(ctx context.Context, number string, i *Invoice) error
	Delete(ctx context.Context, number string) error
}
