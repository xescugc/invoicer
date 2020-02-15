package billing_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/template"
)

func TestCreateTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "can",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		b.Templates.EXPECT().Find(ctx, et.Canonical).Return(nil, billing.ErrNotFoundTemplate)
		b.Templates.EXPECT().Create(ctx, &et).Return(nil)

		err := b.Billing.CreateTemplate(ctx, &et)
		require.NoError(t, err)
	})
	t.Run("ErrAlreadyExistsTemplate", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "can",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		b.Templates.EXPECT().Find(ctx, et.Canonical).Return(&et, nil)

		err := b.Billing.CreateTemplate(ctx, &et)
		assert.EqualError(t, err, billing.ErrAlreadyExistsTemplate.Error())
	})
	t.Run("ErrInvalidTemplateCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		err := b.Billing.CreateTemplate(ctx, &et)
		assert.EqualError(t, err, billing.ErrInvalidTemplateCanonical.Error())
	})
}

func TestGetTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "can",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		b.Templates.EXPECT().Find(ctx, et.Canonical).Return(&et, nil)

		c, err := b.Billing.GetTemplate(ctx, et.Canonical)
		require.NoError(t, err)
		assert.Equal(t, &et, c)
	})
	t.Run("ErrInvalidTemplateCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		c, err := b.Billing.GetTemplate(ctx, et.Canonical)
		assert.EqualError(t, err, billing.ErrInvalidTemplateCanonical.Error())
		assert.Nil(t, c)
	})
}

func TestGetTemplates(t *testing.T) {
	var (
		b   = NewMockBilling(t)
		ctx = context.Background()
		et  = template.Template{
			Name:      "name",
			Canonical: "can",
			Template:  []byte("template"),
		}
	)
	defer b.Finish()

	b.Templates.EXPECT().Filter(ctx).Return([]*template.Template{&et}, nil)

	cs, err := b.Billing.GetTemplates(ctx)
	require.NoError(t, err)
	assert.Equal(t, []*template.Template{&et}, cs)
}

func TestUpdateTemplate(t *testing.T) {
	t.Run("SuccessWithSameCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "can",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		b.Templates.EXPECT().Find(ctx, et.Canonical).Return(&et, nil)
		b.Templates.EXPECT().Update(ctx, et.Canonical, &et).Return(nil)

		err := b.Billing.UpdateTemplate(ctx, et.Canonical, &et)
		require.NoError(t, err)
	})
	t.Run("SuccessWithDiffCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "newcan",
				Template:  []byte("template"),
			}
			oldcan = "can"
		)
		defer b.Finish()

		b.Templates.EXPECT().Find(ctx, oldcan).Return(&et, nil)
		b.Templates.EXPECT().Find(ctx, et.Canonical).Return(nil, billing.ErrNotFoundTemplate)
		b.Templates.EXPECT().Update(ctx, oldcan, &et).Return(nil)

		err := b.Billing.UpdateTemplate(ctx, oldcan, &et)
		require.NoError(t, err)
	})
	t.Run("ErrWithDiffCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "can",
				Template:  []byte("template"),
			}
			diffCan = "diffcan"
		)
		defer b.Finish()

		b.Templates.EXPECT().Find(ctx, et.Canonical).Return(&et, nil)
		b.Templates.EXPECT().Find(ctx, diffCan).Return(&et, nil)

		err := b.Billing.UpdateTemplate(ctx, diffCan, &et)
		assert.EqualError(t, err, billing.ErrAlreadyExistsTemplate.Error())
	})
	t.Run("ErrInvalidTemplateCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		err := b.Billing.UpdateTemplate(ctx, et.Canonical, &et)
		assert.EqualError(t, err, billing.ErrInvalidTemplateCanonical.Error())
	})
	t.Run("ErrInvalidTemplateCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "can",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		err := b.Billing.UpdateTemplate(ctx, "", &et)
		assert.EqualError(t, err, billing.ErrInvalidTemplateCanonical.Error())
	})
}

func TestDeleteTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "can",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		b.Templates.EXPECT().Delete(ctx, et.Canonical).Return(nil)

		err := b.Billing.DeleteTemplate(ctx, et.Canonical)
		require.NoError(t, err)
	})
	t.Run("ErrInvalidTemplateCanonical", func(t *testing.T) {
		var (
			b   = NewMockBilling(t)
			ctx = context.Background()
			et  = template.Template{
				Name:      "name",
				Canonical: "",
				Template:  []byte("template"),
			}
		)
		defer b.Finish()

		err := b.Billing.DeleteTemplate(ctx, et.Canonical)
		assert.EqualError(t, err, billing.ErrInvalidTemplateCanonical.Error())
	})
}
