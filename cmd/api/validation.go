package api

import "packages/validation"

func init() {
	validation.RegisterFieldTrans(map[string]string{
		"Email": "メールアドレス",
	})
}
