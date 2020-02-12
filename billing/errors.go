package billing

import "errors"

var (
	ErrNotFoundUser     = errors.New("user not found")
	ErrNotFoundCustomer = errors.New("customer not found")
	ErrNotFoundInvoice  = errors.New("invoice not found")

	ErrAlreadyExistsUser     = errors.New("user already exists")
	ErrAlreadyExistsCustomer = errors.New("customer already exists")

	ErrInvalidCustomerCanonical = errors.New("invalid customer canonical")
	ErrInvalidInvoiceNumber     = errors.New("invalid invoice number")
)
