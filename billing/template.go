package billing

import (
	"context"
	"errors"
	"fmt"

	"github.com/xescugc/invoicer/template"

	"github.com/gosimple/slug"
)

func (b *billing) CreateTemplate(ctx context.Context, t *template.Template) error {
	if !slug.IsSlug(t.Canonical) {
		return ErrInvalidTemplateCanonical
	}

	nt, err := b.templates.Find(ctx, t.Canonical)
	if err != nil && !errors.Is(err, ErrNotFoundTemplate) {
		return fmt.Errorf("could not get Template: %w", err)
	}

	if nt != nil {
		return ErrAlreadyExistsTemplate
	}

	err = b.templates.Create(ctx, t)
	if err != nil {
		return fmt.Errorf("could not create Template: %w", err)
	}

	return nil
}

func (b *billing) GetTemplate(ctx context.Context, can string) (*template.Template, error) {
	if !slug.IsSlug(can) {
		return nil, ErrInvalidTemplateCanonical
	}

	t, err := b.templates.Find(ctx, can)
	if err != nil {
		return nil, fmt.Errorf("could not get Template: %w", err)
	}

	return t, nil
}

func (b *billing) GetTemplates(ctx context.Context) ([]*template.Template, error) {
	ts, err := b.templates.Filter(ctx)
	if err != nil {
		return nil, fmt.Errorf("could not get Templates: %w", err)
	}

	return ts, nil
}

func (b *billing) UpdateTemplate(ctx context.Context, can string, t *template.Template) error {
	if !slug.IsSlug(can) {
		return ErrInvalidTemplateCanonical
	}
	if !slug.IsSlug(t.Canonical) {
		return ErrInvalidTemplateCanonical
	}

	// Check that the old 'can' exists
	_, err := b.templates.Find(ctx, can)
	if err != nil {
		return fmt.Errorf("could not get Template: %w", err)
	}

	if can != t.Canonical {
		// Check that the new canonical does not exits
		nt, err := b.templates.Find(ctx, t.Canonical)
		if err != nil && !errors.Is(err, ErrNotFoundTemplate) {
			return fmt.Errorf("could not get Template: %w", err)
		}
		if nt != nil {
			return ErrAlreadyExistsTemplate
		}
	}

	err = b.templates.Update(ctx, can, t)
	if err != nil {
		return fmt.Errorf("could not update Template: %w", err)
	}

	return nil
}

func (b *billing) DeleteTemplate(ctx context.Context, can string) error {
	if !slug.IsSlug(can) {
		return ErrInvalidTemplateCanonical
	}

	err := b.templates.Delete(ctx, can)
	if err != nil {
		return fmt.Errorf("could not delete Template: %w", err)
	}

	return nil
}
