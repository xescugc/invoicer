package filesystem

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/customer"
)

type CustomerRepository struct {
	baseDir   string
	entityDir string
}

func NewCustomerRepository(baseDir string) (*CustomerRepository, error) {
	cr := &CustomerRepository{
		baseDir:   baseDir,
		entityDir: filepath.Join(baseDir, "customers"),
	}

	err := os.Mkdir(cr.entityDir, os.ModeDir|0700)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}

	return cr, nil
}

func (r *CustomerRepository) Create(ctx context.Context, c *customer.Customer) error {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.entityDir, jsonFilename(c.Canonical)), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *CustomerRepository) Find(ctx context.Context, canonical string) (*customer.Customer, error) {
	c, err := r.find(ctx, filepath.Join(r.entityDir, jsonFilename(canonical)))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (r *CustomerRepository) find(ctx context.Context, path string) (*customer.Customer, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, billing.ErrNotFoundCustomer
		}
		return nil, err
	}

	var c customer.Customer
	err = json.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (r *CustomerRepository) Filter(ctx context.Context) ([]*customer.Customer, error) {
	files, err := ioutil.ReadDir(r.entityDir)
	if err != nil {
		return nil, err
	}

	cs := make([]*customer.Customer, 0, len(files))
	for _, f := range files {
		c, err := r.find(ctx, filepath.Join(r.entityDir, f.Name()))
		if err != nil {
			return nil, err
		}

		cs = append(cs, c)
	}
	return cs, nil
}

func (r *CustomerRepository) Update(ctx context.Context, can string, c *customer.Customer) error {
	b, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	if can != c.Canonical {
		err := r.Delete(ctx, can)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(filepath.Join(r.entityDir, jsonFilename(c.Canonical)), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *CustomerRepository) Delete(ctx context.Context, can string) error {
	err := os.Remove(filepath.Join(r.entityDir, jsonFilename(can)))
	if err != nil {
		return err
	}

	return nil
}
