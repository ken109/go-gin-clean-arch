package driver

import (
	"fmt"

	"go-gin-clean-arch/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewRDB() (*gorm.DB, error) {
	var db *gorm.DB

	var con string

	if config.Env.DB.Socket != "" {
		con = fmt.Sprintf("unix(%s)", config.Env.DB.Socket)
	} else {
		con = fmt.Sprintf("tcp(%s:%d)", config.Env.DB.Host, config.Env.DB.Port)
	}

	dsn := fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Env.DB.User,
		config.Env.DB.Password,
		con,
		config.Env.DB.Name,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
