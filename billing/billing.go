package billing

import (
	"context"

	"github.com/gosimple/slug"
	"github.com/xescugc/invoicer/customer"
	"github.com/xescugc/invoicer/invoice"
	"github.com/xescugc/invoicer/user"
)

type Billing interface {
	CreateUser(ctx context.Context, u *user.User) error
	GetUser(ctx context.Context) (*user.User, error)
	UpdateUser(ctx context.Context, u *user.User) error

	CreateCustomer(ctx context.Context, c *customer.Customer) error
	GetCustomers(ctx context.Context) ([]*customer.Customer, error)
	GetCustomer(ctx context.Context, canonical string) (*customer.Customer, error)
	UpdateCustomer(ctx context.Context, canonical string, c *customer.Customer) error
	DeleteCustomer(ctx context.Context, canonical string) error

	CreateInvoice(ctx context.Context, c *invoice.Invoice, cusCan string) error
	GetInvoices(ctx context.Context) ([]*invoice.Invoice, error)
	GetInvoice(ctx context.Context, number string) (*invoice.Invoice, error)
	UpdateInvoice(ctx context.Context, number string, c *invoice.Invoice, cusCan string) error
	DeleteInvoice(ctx context.Context, number string) error
}

type billing struct {
	users     user.Repository
	customers customer.Repository
	invoices  invoice.Repository
}

func New(ur user.Repository, cr customer.Repository, ir invoice.Repository) Billing {
	return &billing{
		users:     ur,
		customers: cr,
		invoices:  ir,
	}
}

func init() {
	slug.MaxLength = 30
}
