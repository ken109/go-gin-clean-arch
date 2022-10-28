package mysql

import (
	"go-gin-ddd/domain"
	"go-gin-ddd/driver"
)

func init() {
	err := driver.GetRDB().AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		panic(err)
	}
}
