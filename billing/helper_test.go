package billing_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/mock"
)

type MockBilling struct {
	Users     *mock.UserRepository
	Customers *mock.CustomerRepository
	Invoices  *mock.InvoiceRepository
	Templates *mock.TemplateRepository

	Ctrl    *gomock.Controller
	Billing billing.Billing
}

func NewMockBilling(t *testing.T) MockBilling {
	var (
		ctrl = gomock.NewController(t)
		ur   = mock.NewUserRepository(ctrl)
		cr   = mock.NewCustomerRepository(ctrl)
		ir   = mock.NewInvoiceRepository(ctrl)
		tr   = mock.NewTemplateRepository(ctrl)
		b    = billing.New(ur, cr, ir, tr)
	)

	return MockBilling{
		Users:     ur,
		Customers: cr,
		Invoices:  ir,
		Templates: tr,

		Ctrl: ctrl,

		Billing: b,
	}
}

func (m *MockBilling) Finish() {
	m.Ctrl.Finish()
}
