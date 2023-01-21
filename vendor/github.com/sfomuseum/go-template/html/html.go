// Package html provides methods for loading HTML (.html) templates with default functions
package html

import (
	"context"
	"fmt"
	"github.com/sfomuseum/go-template/funcs"
	"html/template"
	"io/fs"
)

// LoadTemplates loads HTML (.html) from 't_fs' with default functions assigned.
func LoadTemplates(ctx context.Context, t_fs ...fs.FS) (*template.Template, error) {

	funcs := TemplatesFuncMap()
	t := template.New("html").Funcs(funcs)

	var err error

	for idx, f := range t_fs {

		t, err = t.ParseFS(f, "*.html")

		if err != nil {
			return nil, fmt.Errorf("Failed to load templates from FS at offset %d, %w", idx, err)
		}
	}

	return t, nil
}

// TemplatesFuncMap() returns a `template.FuncMap` instance with default functions assigned.
func TemplatesFuncMap() template.FuncMap {

	return template.FuncMap{
		// For example: {{ if (IsAvailable "Account" .) }}
		"IsAvailable": funcs.IsAvailable,
		"Add":         funcs.Add,
		"JoinPath": funcs.JoinPath,
	}
}
