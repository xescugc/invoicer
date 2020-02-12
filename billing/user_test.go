package billing_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/user"
)

func TestCreateUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			eu  = user.User{
				Name:      "name",
				Address:   "address",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		b.Users.EXPECT().Find(ctx).Return(nil, billing.ErrNotFoundUser)
		b.Users.EXPECT().Create(ctx, &eu).Return(nil)

		err := b.Billing.CreateUser(ctx, &eu)
		require.NoError(t, err)
	})
	t.Run("ErrAlreadyExistsUser", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			eu  = user.User{
				Name:      "name",
				Address:   "address",
				VATNumber: "vatnumber",
			}
		)
		defer b.Finish()

		b.Users.EXPECT().Find(ctx).Return(&eu, nil)

		err := b.Billing.CreateUser(ctx, &eu)
		assert.EqualError(t, err, billing.ErrAlreadyExistsUser.Error())
	})
}

func TestGetUser(t *testing.T) {
	var (
		b   = NewMockBilling(t)
		ctx = context.Background()
		eu  = user.User{
			Name:      "name",
			Address:   "address",
			VATNumber: "vatnumber",
		}
	)
	defer b.Finish()

	b.Users.EXPECT().Find(ctx).Return(&eu, nil)

	u, err := b.Billing.GetUser(ctx)
	require.NoError(t, err)
	assert.Equal(t, &eu, u)
}

func TestUpdateUser(t *testing.T) {
	var (
		b   = NewMockBilling(t)
		ctx = context.Background()
		eu  = user.User{
			Name:      "name",
			Address:   "address",
			VATNumber: "vatnumber",
		}
	)
	defer b.Finish()

	b.Users.EXPECT().Update(ctx, &eu).Return(nil)

	err := b.Billing.UpdateUser(ctx, &eu)
	require.NoError(t, err)
}
