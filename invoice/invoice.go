package invoice

import (
	"fmt"
	"time"

	"github.com/xescugc/invoicer/customer"
	"github.com/xescugc/invoicer/user"
	"golang.org/x/text/currency"
)

type Invoice struct {
	Number string

	Date time.Time

	User     user.User
	Customer customer.Customer

	Items []Item
	VAT   float64
}

func (i Invoice) Total() string {
	total := 0.0
	curr := currency.XXX
	for _, it := range i.Items {
		total += it.Price
		curr = it.Currency
	}

	return fmt.Sprintf("%g %s", total, curr)
}
