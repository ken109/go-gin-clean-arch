package api

import "go-gin-ddd/packages/validation"

func init() {
	validation.RegisterFieldTrans(map[string]string{
		"Email": "メールアドレス",
	})
}
