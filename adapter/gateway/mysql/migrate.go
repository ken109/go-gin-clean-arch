package mysql

import (
	"go-gin-clean-arch/domain"
	"go-gin-clean-arch/driver"
)

func init() {
	err := driver.GetRDB().AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		panic(err)
	}
}
