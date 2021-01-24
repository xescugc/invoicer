package invoice

import (
	"fmt"
	"time"

	"github.com/xescugc/invoicer/customer"
	"github.com/xescugc/invoicer/user"
	"github.com/xescugc/marshaler"
)

type Invoice struct {
	Number string

	Date time.Time

	User     user.User
	Customer customer.Customer

	Items []Item
	VAT   float64

	Currency marshaler.CurrencyUnit
}

func (i Invoice) Total() string {
	total := 0.0
	for _, it := range i.Items {
		total += it.Price
	}

	return fmt.Sprintf("%g %s", total, i.Currency)
}
