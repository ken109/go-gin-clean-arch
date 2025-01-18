package mysql

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"go-gin-clean-arch/config"
	"go-gin-clean-arch/packages/cerrors"
)

func getDB(ctx context.Context) *gorm.DB {
	db, _ := ctx.Value(config.DBKey).(*gorm.DB)
	return db
}

func dbError(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return cerrors.NotFound()
	default:
		return cerrors.NewUnexpected(err)
	}
}
