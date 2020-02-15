package billing_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/mock"
)

func TestNew(t *testing.T) {
	var (
		ctrl = gomock.NewController(t)
		ur   = mock.NewUserRepository(ctrl)
		cr   = mock.NewCustomerRepository(ctrl)
		ir   = mock.NewInvoiceRepository(ctrl)
		tr   = mock.NewTemplateRepository(ctrl)
	)

	b := billing.New(ur, cr, ir, tr)
	assert.NotNil(t, b)

	defer ctrl.Finish()
}
