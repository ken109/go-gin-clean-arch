package mysql

import (
	"go-gin-clean-arch/packages/errors"
	"gorm.io/gorm"
)

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
