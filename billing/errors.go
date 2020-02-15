package billing

import "errors"

var (
	ErrNotFoundUser     = errors.New("user not found")
	ErrNotFoundCustomer = errors.New("customer not found")
	ErrNotFoundInvoice  = errors.New("invoice not found")
	ErrNotFoundTemplate = errors.New("template not found")

	ErrAlreadyExistsUser     = errors.New("user already exists")
	ErrAlreadyExistsCustomer = errors.New("customer already exists")
	ErrAlreadyExistsTemplate = errors.New("template already exists")

	ErrInvalidCustomerCanonical = errors.New("invalid customer canonical")
	ErrInvalidTemplateCanonical = errors.New("invalid template canonical")
	ErrInvalidInvoiceNumber     = errors.New("invalid invoice number")
)
