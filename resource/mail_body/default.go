package mail_body

import (
	_ "embed"
	"html/template"
)

var defaultHTMLTemplateString = `<!DOCTYPE html>
<div>
    <h1>{{.Title}}</h1>
    <p>{{.Body}}</p>
</div>
`
var defaultHTMLTemplate *template.Template

func init() {
	var err error
	defaultHTMLTemplate, err = template.New("defaultHtml").Parse(defaultHTMLTemplateString)
	if err != nil {
		panic(err)
	}
}

type Default struct {
	Title string
	Body  string
}

func (d Default) Subject() string {
	return d.Title
}

func (d Default) HTML() (string, error) {
	return setHTMLTemplate(defaultHTMLTemplate, d)
}

func (d Default) Plain() (string, error) {
	return htmlToPlain(d)
}
