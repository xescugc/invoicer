package billing

import (
	"bytes"
	"context"
	"fmt"
	"text/template"

	"github.com/gosimple/slug"
	"github.com/xescugc/invoicer/invoice"
)

func (b *billing) CreateInvoice(ctx context.Context, i *invoice.Invoice, cusCan string) error {
	if !slug.IsSlug(cusCan) {
		return ErrInvalidCustomerCanonical
	}
	if i.Number == "" {
		return ErrInvalidInvoiceNumber
	}

	c, err := b.customers.Find(ctx, cusCan)
	if err != nil {
		return fmt.Errorf("could not get Customer: %w", err)
	}

	// TODO Check number and time
	u, err := b.users.Find(ctx)
	if err != nil {
		return err
	}

	i.User = *u
	i.Customer = *c

	err = b.invoices.Create(ctx, i)
	if err != nil {
		return err
	}

	return nil
}

func (b *billing) GetInvoices(ctx context.Context) ([]*invoice.Invoice, error) {
	is, err := b.invoices.Filter(ctx)
	if err != nil {
		return nil, err
	}

	return is, nil
}

func (b *billing) GetInvoice(ctx context.Context, number string) (*invoice.Invoice, error) {
	if !isValidInvoiceNumber(number) {
		return nil, ErrInvalidInvoiceNumber
	}

	i, err := b.invoices.Find(ctx, number)
	if err != nil {
		return nil, err
	}

	return i, nil
}

func (b *billing) UpdateInvoice(ctx context.Context, number string, i *invoice.Invoice, cusCan string) error {
	if !isValidInvoiceNumber(number) {
		return ErrInvalidInvoiceNumber
	}

	dbin, err := b.GetInvoice(ctx, number)
	if err != nil {
		return err
	}

	if cusCan == "" {
		cusCan = dbin.Customer.Canonical
	}

	// TODO check number like on the Customer
	if dbin.Customer.Canonical != cusCan {
		c, err := b.GetCustomer(ctx, cusCan)
		if err != nil {
			return err
		}

		i.Customer = *c
	} else {
		i.Customer = dbin.Customer
	}

	i.User = dbin.User

	err = b.invoices.Update(ctx, number, i)
	if err != nil {
		return err
	}

	return nil
}

func (b *billing) DeleteInvoice(ctx context.Context, number string) error {
	if !isValidInvoiceNumber(number) {
		return ErrInvalidInvoiceNumber
	}

	err := b.invoices.Delete(ctx, number)
	if err != nil {
		return err
	}

	return nil
}

func (b *billing) ViewInvoice(ctx context.Context, invoiceNumber, templateCan string) ([]byte, error) {
	if !isValidInvoiceNumber(invoiceNumber) {
		return nil, ErrInvalidInvoiceNumber
	}

	in, err := b.invoices.Find(ctx, invoiceNumber)
	if err != nil {
		return nil, err
	}

	tpl, err := b.templates.Find(ctx, templateCan)
	if err != nil {
		return nil, fmt.Errorf("could not get Templates: %w", err)
	}

	t, err := template.New("invoice.html").Parse(string(tpl.Template))
	if err != nil {
		return nil, err
	}

	buff := bytes.Buffer{}

	err = t.Execute(&buff, in)
	if err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
