package mysql

import (
	"go-gin-ddd/domain"
	"go-gin-ddd/driver/rdb"
)

func init() {
	err := rdb.Get().AutoMigrate(
		&domain.User{},
	)
	if err != nil {
		panic(err)
	}
}
