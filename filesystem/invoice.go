package filesystem

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/invoice"
)

type InvoiceRepository struct {
	baseDir   string
	entityDir string
}

func NewInvoiceRepository(baseDir string) (*InvoiceRepository, error) {
	ir := &InvoiceRepository{
		baseDir:   baseDir,
		entityDir: filepath.Join(baseDir, "invoices"),
	}

	err := os.Mkdir(ir.entityDir, os.ModeDir|0700)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}

	return ir, nil
}

func (r *InvoiceRepository) Create(ctx context.Context, i *invoice.Invoice) error {
	b, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.entityDir, jsonFilename(i.Number)), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *InvoiceRepository) Find(ctx context.Context, number string) (*invoice.Invoice, error) {
	c, err := r.find(ctx, filepath.Join(r.entityDir, jsonFilename(number)))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *InvoiceRepository) find(ctx context.Context, path string) (*invoice.Invoice, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, billing.ErrNotFoundInvoice
		}
		return nil, err
	}

	var i invoice.Invoice
	err = json.Unmarshal(b, &i)
	if err != nil {
		return nil, err
	}

	return &i, nil
}

func (r *InvoiceRepository) Filter(ctx context.Context) ([]*invoice.Invoice, error) {
	files, err := ioutil.ReadDir(r.entityDir)
	if err != nil {
		return nil, err
	}

	is := make([]*invoice.Invoice, 0, len(files))
	for _, f := range files {
		i, err := r.find(ctx, filepath.Join(r.entityDir, f.Name()))
		if err != nil {
			return nil, err
		}

		is = append(is, i)
	}
	return is, nil
}

func (r *InvoiceRepository) Update(ctx context.Context, number string, i *invoice.Invoice) error {
	b, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		return err
	}

	if number != i.Number {
		err := r.Delete(ctx, number)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(filepath.Join(r.entityDir, jsonFilename(i.Number)), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *InvoiceRepository) Delete(ctx context.Context, number string) error {
	err := os.Remove(filepath.Join(r.entityDir, jsonFilename(number)))
	if err != nil {
		return err
	}

	return nil
}
