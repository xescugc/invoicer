package template

import "context"

//go:generate mockgen -destination=../mock/template_repository.go -mock_names=Repository=TemplateRepository -package mock github.com/xescugc/invoicer/template Repository

type Repository interface {
	Create(ctx context.Context, t *Template) error
	Find(ctx context.Context, can string) (*Template, error)
	Filter(ctx context.Context) ([]*Template, error)
	Update(ctx context.Context, can string, t *Template) error
	Delete(ctx context.Context, can string) error
}
