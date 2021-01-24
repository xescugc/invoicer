package user

import "github.com/xescugc/invoicer/address"

type User struct {
	Name      string
	Address   address.Address
	VATNumber string
}
