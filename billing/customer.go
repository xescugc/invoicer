package billing

import (
	"context"
	"errors"
	"fmt"

	"github.com/xescugc/invoicer/customer"

	"github.com/gosimple/slug"
)

func (b *billing) CreateCustomer(ctx context.Context, c *customer.Customer) error {
	if !slug.IsSlug(c.Canonical) {
		return ErrInvalidCustomerCanonical
	}

	nc, err := b.customers.Find(ctx, c.Canonical)
	if err != nil && !errors.Is(err, ErrNotFoundCustomer) {
		return fmt.Errorf("could not get Customer: %w", err)
	}

	if nc != nil {
		return ErrAlreadyExistsCustomer
	}

	err = b.customers.Create(ctx, c)
	if err != nil {
		return fmt.Errorf("could not create Customer: %w", err)
	}

	return nil
}

func (b *billing) GetCustomer(ctx context.Context, can string) (*customer.Customer, error) {
	if !slug.IsSlug(can) {
		return nil, ErrInvalidCustomerCanonical
	}

	c, err := b.customers.Find(ctx, can)
	if err != nil {
		return nil, fmt.Errorf("could not get Customer: %w", err)
	}

	return c, nil
}

func (b *billing) GetCustomers(ctx context.Context) ([]*customer.Customer, error) {
	cs, err := b.customers.Filter(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get Customers: %w", err)
	}

	return cs, nil
}

func (b *billing) UpdateCustomer(ctx context.Context, can string, c *customer.Customer) error {
	if !slug.IsSlug(can) {
		return ErrInvalidCustomerCanonical
	}
	if !slug.IsSlug(c.Canonical) {
		return ErrInvalidCustomerCanonical
	}

	// Check that the old 'can' exists
	_, err := b.customers.Find(ctx, can)
	if err != nil {
		return fmt.Errorf("could not get Customer: %w", err)
	}

	if can != c.Canonical {
		// Check that the new canonical does not exits
		nc, err := b.customers.Find(ctx, c.Canonical)
		if err != nil && !errors.Is(err, ErrNotFoundCustomer) {
			return fmt.Errorf("could not get Customer: %w", err)
		}
		if nc != nil {
			return ErrAlreadyExistsCustomer
		}
	}

	err = b.customers.Update(ctx, can, c)
	if err != nil {
		return fmt.Errorf("could not update Customer: %w", err)
	}

	return nil
}

func (b *billing) DeleteCustomer(ctx context.Context, can string) error {
	if !slug.IsSlug(can) {
		return ErrInvalidCustomerCanonical
	}

	err := b.customers.Delete(ctx, can)
	if err != nil {
		return fmt.Errorf("could not delete Customer: %w", err)
	}

	return nil
}
