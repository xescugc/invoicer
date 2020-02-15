package filesystem

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/xescugc/invoicer/billing"
	"github.com/xescugc/invoicer/template"
)

type TemplateRepository struct {
	baseDir   string
	entityDir string
}

func NewTemplateRepository(baseDir string) (*TemplateRepository, error) {
	tr := &TemplateRepository{
		baseDir:   baseDir,
		entityDir: filepath.Join(baseDir, "templates"),
	}

	err := os.Mkdir(tr.entityDir, os.ModeDir|0700)
	if err != nil && !errors.Is(err, os.ErrExist) {
		return nil, err
	}

	return tr, nil
}

func (r *TemplateRepository) Create(ctx context.Context, t *template.Template) error {
	b, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(r.entityDir, jsonFilename(t.Canonical)), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *TemplateRepository) Find(ctx context.Context, canonical string) (*template.Template, error) {
	t, err := r.find(ctx, filepath.Join(r.entityDir, jsonFilename(canonical)))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func (r *TemplateRepository) find(ctx context.Context, path string) (*template.Template, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, billing.ErrNotFoundTemplate
		}
		return nil, err
	}

	var t template.Template
	err = json.Unmarshal(b, &t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (r *TemplateRepository) Filter(ctx context.Context) ([]*template.Template, error) {
	files, err := ioutil.ReadDir(r.entityDir)
	if err != nil {
		return nil, err
	}

	ts := make([]*template.Template, 0, len(files))
	for _, f := range files {
		t, err := r.find(ctx, filepath.Join(r.entityDir, f.Name()))
		if err != nil {
			return nil, err
		}

		ts = append(ts, t)
	}
	return ts, nil
}

func (r *TemplateRepository) Update(ctx context.Context, can string, t *template.Template) error {
	b, err := json.MarshalIndent(t, "", " ")
	if err != nil {
		return err
	}

	if can != t.Canonical {
		err := r.Delete(ctx, can)
		if err != nil {
			return err
		}
	}

	err = ioutil.WriteFile(filepath.Join(r.entityDir, jsonFilename(t.Canonical)), b, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (r *TemplateRepository) Delete(ctx context.Context, can string) error {
	err := os.Remove(filepath.Join(r.entityDir, jsonFilename(can)))
	if err != nil {
		return err
	}

	return nil
}
