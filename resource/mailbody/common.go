package mailbody

import (
	"bytes"
	templateHtml "html/template"

	"jaytaylor.com/html2text"

	"go-gin-clean-arch/packages/errors"
)

type MailBody interface {
	Subject() string
	HTML() (string, error)
	Plain() (string, error)
}

func htmlToPlain(body MailBody) (string, error) {
	html, err := body.HTML()
	if err != nil {
		return "", err
	}
	return html2text.FromString(html)
}

func setHTMLTemplate(template *templateHtml.Template, data interface{}) (string, error) {
	var out bytes.Buffer
	err := template.Execute(&out, data)
	if err != nil {
		return "", errors.NewUnexpected(err)
	}
	return out.String(), nil
}

// func setPlainTemplate(template *templateText.Template, data interface{}) (string, error) {
// 	var out bytes.Buffer
// 	err := template.Execute(&out, data)
// 	if err != nil {
// 		return "", errors.NewUnexpected(err)
// 	}
// 	return out.String(), nil
// }
