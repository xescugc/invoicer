package billing_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/customer"
)

func TestCreateCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		b.Customers.EXPECT().Find(ctx, ec.Canonical).Return(nil, billing.ErrNotFoundCustomer)
		b.Customers.EXPECT().Create(ctx, &ec).Return(nil)

		err := b.Billing.CreateCustomer(ctx, &ec)
		require.NoError(t, err)
	})
	t.Run("ErrAlreadyExistsCustomer", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		b.Customers.EXPECT().Find(ctx, ec.Canonical).Return(&ec, nil)

		err := b.Billing.CreateCustomer(ctx, &ec)
		assert.EqualError(t, err, billing.ErrAlreadyExistsCustomer.Error())
	})
	t.Run("ErrInvalidCustomerCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		err := b.Billing.CreateCustomer(ctx, &ec)
		assert.EqualError(t, err, billing.ErrInvalidCustomerCanonical.Error())
	})
}

func TestGetCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		b.Customers.EXPECT().Find(ctx, ec.Canonical).Return(&ec, nil)

		c, err := b.Billing.GetCustomer(ctx, ec.Canonical)
		require.NoError(t, err)
		assert.Equal(t, &ec, c)
	})
	t.Run("ErrInvalidCustomerCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		c, err := b.Billing.GetCustomer(ctx, ec.Canonical)
		assert.EqualError(t, err, billing.ErrInvalidCustomerCanonical.Error())
		assert.Nil(t, c)
	})
}

func TestGetCustomers(t *testing.T) {
	var (
		b   = NewMockBilling(t)
		ctx = context.Background()
		ec  = customer.Customer{
			Name:      "name",
			Canonical: "can",
			VATNumber: "vatnumber",
		}
	)
	defer b.Finish()

	b.Customers.EXPECT().Filter(ctx).Return([]*customer.Customer{&ec}, nil)

	cs, err := b.Billing.GetCustomers(ctx)
	require.NoError(t, err)
	assert.Equal(t, []*customer.Customer{&ec}, cs)
}

func TestUpdateCustomer(t *testing.T) {
	t.Run("SuccessWithSameCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		b.Customers.EXPECT().Find(ctx, ec.Canonical).Return(&ec, nil)
		b.Customers.EXPECT().Update(ctx, ec.Canonical, &ec).Return(nil)

		err := b.Billing.UpdateCustomer(ctx, ec.Canonical, &ec)
		require.NoError(t, err)
	})
	t.Run("SuccessWithDiffCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "newcan",
				VATNumber: "vatnumber",
			}
			oldcan = "can"
		)
		defer b.Finish()

		b.Customers.EXPECT().Find(ctx, oldcan).Return(&ec, nil)
		b.Customers.EXPECT().Find(ctx, ec.Canonical).Return(nil, billing.ErrNotFoundCustomer)
		b.Customers.EXPECT().Update(ctx, oldcan, &ec).Return(nil)

		err := b.Billing.UpdateCustomer(ctx, oldcan, &ec)
		require.NoError(t, err)
	})
	t.Run("ErrWithDiffCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
			diffCan = "diffcan"
		)
		defer b.Finish()

		b.Customers.EXPECT().Find(ctx, ec.Canonical).Return(&ec, nil)
		b.Customers.EXPECT().Find(ctx, diffCan).Return(&ec, nil)

		err := b.Billing.UpdateCustomer(ctx, diffCan, &ec)
		assert.EqualError(t, err, billing.ErrAlreadyExistsCustomer.Error())
	})
	t.Run("ErrInvalidCustomerCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		err := b.Billing.UpdateCustomer(ctx, ec.Canonical, &ec)
		assert.EqualError(t, err, billing.ErrInvalidCustomerCanonical.Error())
	})
	t.Run("ErrInvalidCustomerCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		err := b.Billing.UpdateCustomer(ctx, "", &ec)
		assert.EqualError(t, err, billing.ErrInvalidCustomerCanonical.Error())
	})
}

func TestDeleteCustomer(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "can",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		b.Customers.EXPECT().Delete(ctx, ec.Canonical).Return(nil)

		err := b.Billing.DeleteCustomer(ctx, ec.Canonical)
		require.NoError(t, err)
	})
	t.Run("ErrInvalidCustomerCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			ec  = customer.Customer{
				Name:      "name",
				Canonical: "",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		err := b.Billing.DeleteCustomer(ctx, ec.Canonical)
		assert.EqualError(t, err, billing.ErrInvalidCustomerCanonical.Error())
	})
}
