package model

import (
	"github.com/xescugc/invoicer/template"
)

type Template struct {
	Name      string
	Canonical string
}

func NewTemplateFromDomain(t *template.Template) Template {
	return Template{
		Name:      t.Name,
		Canonical: t.Canonical,
	}
}

func (t Template) ToDomain() *template.Template {
	return &template.Template{
		Name:      t.Name,
		Canonical: t.Canonical,
	}
}
