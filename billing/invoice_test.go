package billing_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/customer"
	"github.com/xescugc/invoicer/invoice"
	"github.com/xescugc/invoicer/template"
	"github.com/xescugc/invoicer/user"
	"golang.org/x/text/currency"
)

func TestCreateInvoice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
			eu = user.User{
				Name: "Pepito",
			}
			in = invoice.Invoice{
				Items: []invoice.Item{
					invoice.Item{
						Description: "Some charge",
						Price:       30,
						Currency:    currency.EUR,
					},
				},
			}
			ein = invoice.Invoice{
				Customer: ec,
				User:     eu,
				Items:    in.Items,
			}
		)
		defer b.Finish()

		b.Customers.EXPECT().Find(ctx, ec.Canonical).Return(&ec, nil)
		b.Users.EXPECT().Find(ctx).Return(&eu, nil)
		b.Invoices.EXPECT().Create(ctx, &ein).Return(nil)

		err := b.Billing.CreateInvoice(ctx, &in, ec.Canonical)
		require.NoError(t, err)
	})
}

func TestGetInvoice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			in  = invoice.Invoice{
				Items: []invoice.Item{
					invoice.Item{
						Description: "Some charge",
						Price:       30,
						Currency:    currency.EUR,
					},
				},
			}
			inum = "number"
		)
		defer b.Finish()

		b.Invoices.EXPECT().Find(ctx, inum).Return(&in, nil)

		i, err := b.Billing.GetInvoice(ctx, inum)
		require.NoError(t, err)
		assert.Equal(t, &in, i)
	})
	t.Run("ErrInvalidInvoiceNumber", func(t *testing.T) {
		var (
			b    = NewMockBilling(t)
			ctx  = context.Background()
			inum = ""
		)
		defer b.Finish()

		i, err := b.Billing.GetInvoice(ctx, inum)
		assert.Nil(t, i)
		assert.EqualError(t, err, billing.ErrInvalidInvoiceNumber.Error())
	})
}

func TestGetInvoices(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			in  = invoice.Invoice{
				Items: []invoice.Item{
					invoice.Item{
						Description: "Some charge",
						Price:       30,
						Currency:    currency.EUR,
					},
				},
			}
		)
		defer b.Finish()

		b.Invoices.EXPECT().Filter(ctx).Return([]*invoice.Invoice{&in}, nil)

		ins, err := b.Billing.GetInvoices(ctx)
		require.NoError(t, err)
		assert.Equal(t, []*invoice.Invoice{&in}, ins)
	})
}

func TestUpdateInvoice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
			in = invoice.Invoice{
				Customer: ec,
				Items: []invoice.Item{
					invoice.Item{
						Description: "Some charge",
						Price:       30,
						Currency:    currency.EUR,
					},
				},
			}
			inum = "number"
		)
		defer b.Finish()

		b.Invoices.EXPECT().Find(ctx, inum).Return(&in, nil)
		b.Invoices.EXPECT().Update(ctx, inum, &in).Return(nil)

		err := b.Billing.UpdateInvoice(ctx, inum, &in, ec.Canonical)
		require.NoError(t, err)
	})
	t.Run("SuccessDiffCustomer", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			enc = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
			in = invoice.Invoice{
				Items: []invoice.Item{
					invoice.Item{
						Description: "Some charge",
						Price:       30,
						Currency:    currency.EUR,
					},
				},
			}
			ein = invoice.Invoice{
				Customer: enc,
				Items: []invoice.Item{
					invoice.Item{
						Description: "Some charge",
						Price:       30,
						Currency:    currency.EUR,
					},
				},
			}
			inum      = "number"
			newCusCan = "newcuscan"
		)
		defer b.Finish()

		b.Invoices.EXPECT().Find(ctx, inum).Return(&in, nil)
		b.Customers.EXPECT().Find(ctx, newCusCan).Return(&enc, nil)
		b.Invoices.EXPECT().Update(ctx, inum, &ein).Return(nil)

		err := b.Billing.UpdateInvoice(ctx, inum, &in, newCusCan)
		require.NoError(t, err)
	})
	t.Run("ErrInvalidInvoiceNumber", func(t *testing.T) {
		var (
			b    = NewMockBilling(t)
			ctx  = context.Background()
			inum = ""
		)
		defer b.Finish()

		i, err := b.Billing.GetInvoice(ctx, inum)
		assert.Nil(t, i)
		assert.EqualError(t, err, billing.ErrInvalidInvoiceNumber.Error())
	})
}

func TestDeleteInvoice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b    = NewMockBilling(t)
			ctx  = context.Background()
			inum = "number"
		)
		defer b.Finish()

		b.Invoices.EXPECT().Delete(ctx, inum).Return(nil)

		err := b.Billing.DeleteInvoice(ctx, inum)
		require.NoError(t, err)
	})
	t.Run("ErrInvalidInvoiceNumber", func(t *testing.T) {
		var (
			b    = NewMockBilling(t)
			ctx  = context.Background()
			inum = ""
		)
		defer b.Finish()

		err := b.Billing.DeleteInvoice(ctx, inum)
		assert.EqualError(t, err, billing.ErrInvalidInvoiceNumber.Error())
	})
}
func TestViewInvoice(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b      = NewMockBilling(t)
			ctx    = context.Background()
			inum   = "number"
			tplcan = "tplcan"
			ec     = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
			eu = user.User{
				Name:      "name",
				Address:   "address",
				VATNumber: "vatnumber",
			}
			in = invoice.Invoice{
				Number:   inum,
				Customer: ec,
				User:     eu,
				Items: []invoice.Item{
					invoice.Item{
						Description: "Some charge",
						Price:       30,
						Currency:    currency.EUR,
					},
				},
			}
			tpl = template.Template{
				Name:      "name",
				Canonical: tplcan,
				Template:  []byte("{{ .Number }} {{ .Customer.Name }} {{ .User.Name }}"),
			}
		)
		defer b.Finish()

		b.Invoices.EXPECT().Find(ctx, inum).Return(&in, nil)
		b.Templates.EXPECT().Find(ctx, tplcan).Return(&tpl, nil)

		tplin, err := b.Billing.ViewInvoice(ctx, inum, tplcan)
		require.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%s %s %s", in.Number, in.Customer.Name, in.User.Name), string(tplin))
	})
}
