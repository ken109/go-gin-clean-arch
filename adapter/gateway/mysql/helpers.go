package mysql

import (
	"context"

	"gorm.io/gorm"

	"go-gin-clean-arch/config"
	"go-gin-clean-arch/packages/errors"
)

func getDB(ctx context.Context) *gorm.DB {
	db, _ := ctx.Value(config.DBKey).(*gorm.DB)
	return db
}

func dbError(err error) error {
	switch err {
	case nil:
		return nil
	case gorm.ErrRecordNotFound:
		return errors.NotFound()
	default:
		return errors.NewUnexpected(err)
	}
}
