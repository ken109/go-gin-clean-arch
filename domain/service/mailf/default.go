package mailf

import (
	"bytes"
	_ "embed"
	"html/template"
)

//go:embed templates/default.gohtml
var defaultTemplateString string
var defaultTemplate *template.Template

func init() {
	var err error
	defaultTemplate, err = template.New("default").Funcs(funcMap).Parse(defaultTemplateString)
	if err != nil {
		panic(err)
	}
}

type defaultMail struct {
	to      string
	subject string
	data    Data
}

type Data struct {
	Title string
	Body  string
}

func NewDefaultMail(to string, subject string, data Data) Mail {
	return &defaultMail{
		to:      to,
		subject: subject,
		data:    data,
	}
}

func (m *defaultMail) To() (string, error) {
	return m.to, nil
}

func (m *defaultMail) Title() (string, error) {
	return m.subject, nil
}

func (m *defaultMail) Content() (string, error) {
	var html bytes.Buffer

	err := defaultTemplate.Execute(&html, m.data)
	if err != nil {
		return "", err
	}

	return html.String(), nil
}
