package mailf

import (
	"encoding/json"
	"html/template"
	"strings"

	"gorm.io/datatypes"
)

type Sender interface {
	Send(mails Mail) error
	BulkSend(mails []Mail) error
}

type Mail interface {
	To() (string, error)
	Title() (string, error)
	Content() (string, error)
}

var funcMap = template.FuncMap{
	"json": func(a datatypes.JSON, key string) interface{} {
		m := map[string]interface{}{}
		if err := json.Unmarshal(a, &m); err != nil {
			panic(err)
		}
		var value interface{}
		keys := strings.Split(key, ".")
		for _, k := range keys {
			if v, ok := m[k].(map[string]interface{}); ok {
				m = v
			} else {
				value = m[k]
			}
		}
		return value
	},
}
