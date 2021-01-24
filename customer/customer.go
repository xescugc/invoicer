package customer

import "github.com/xescugc/invoicer/address"

type Customer struct {
	Name      string
	Canonical string
	Address   address.Address
	VATNumber string
}
