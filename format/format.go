package format

import (
	"bytes"
	"text/template"
)

type Template struct {
	template *template.Template
}

func NewTemplate(name string, text string) (*Template, error) {

	tmpl, err := template.New(name).Parse(text)
	if err != nil {
		return nil, err
	}

	return &Template{tmpl}, nil
}

func (f *Template) Apply(data map[string]string) (string, error) {

	buf := new(bytes.Buffer)

	if err := f.template.Execute(buf, data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
