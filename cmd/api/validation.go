package main

import "go-gin-clean-arch/packages/validation"

func init() {
	validation.RegisterFieldTrans(map[string]string{
		"Email": "メールアドレス",
	})
}
