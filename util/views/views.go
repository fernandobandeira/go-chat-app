package views

import (
	"fmt"
	"html/template"
	"path/filepath"
)

type Base struct {
	path         string
	baseTemplate *template.Template
}

func NewBase(path string, sharedFolder string) *Base {
	if path[len(path)-1] == '/' {
		path = path[:len(path)-1]
	}
	templateGlob := filepath.FromSlash(fmt.Sprintf("%s/%s/**/*.html", path, sharedFolder))
	t, err := template.ParseGlob(templateGlob)
	if err != nil {
		panic(err)
	}

	v := &Base{
		path:         path,
		baseTemplate: t,
	}

	return v
}

func (v *Base) Parse(view string) *template.Template {
	viewPath := filepath.FromSlash(filepath.Join(v.path, view))
	t, err := v.baseTemplate.Clone()
	if err != nil {
		panic(err)
	}
	t, err = t.ParseFiles(viewPath)
	if err != nil {
		panic(err)
	}
	return t
}
