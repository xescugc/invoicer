package model

import (
	"time"

	"github.com/xescugc/invoicer/invoice"
	"golang.org/x/text/currency"
)

const (
	DefaultDateFormat = "02/01/2006"
)

type Invoice struct {
	Number string

	Date string

	Items []Item

	VAT float64
}

type Item struct {
	Description string
	Price       float64
	Currency    string
}

func NewInvoice() Invoice {
	d := time.Now().Format(DefaultDateFormat)
	return Invoice{
		Items: make([]Item, 1),
		Date:  d,
	}
}

func NewInvoiceFromDomain(i *invoice.Invoice) Invoice {
	items := make([]Item, 0, len(i.Items))

	for _, it := range i.Items {
		itm := Item{
			Description: it.Description,
			Price:       it.Price,
			Currency:    it.Currency.String(),
		}

		items = append(items, itm)
	}

	return Invoice{
		Number: i.Number,
		VAT:    i.VAT,
		Items:  items,
		Date:   i.Date.Format(DefaultDateFormat),
	}
}

func (i Invoice) ToDomain() (*invoice.Invoice, error) {
	items := make([]invoice.Item, 0, len(i.Items))

	for _, it := range i.Items {
		cr, err := currency.ParseISO(it.Currency)
		if err != nil {
			return nil, err
		}
		itm := invoice.Item{
			Description: it.Description,
			Price:       it.Price,
			Currency:    cr,
		}

		items = append(items, itm)
	}

	d, err := time.Parse(DefaultDateFormat, i.Date)
	if err != nil {
		return nil, err
	}

	return &invoice.Invoice{
		Number: i.Number,
		VAT:    i.VAT,
		Items:  items,
		Date:   d,
	}, nil
}
