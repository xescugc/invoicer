package invoice

import "golang.org/x/text/currency"

type Item struct {
	Description string
	Price       float64
	Currency    currency.Unit
}
