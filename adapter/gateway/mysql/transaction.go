package mysql

import (
	"context"
	"go-gin-clean-arch/config"
	"go-gin-clean-arch/usecase"
	"gorm.io/gorm"
)

type transaction struct{}

func NewTransaction() usecase.TransactionRepository {
	return &transaction{}
}

func (tx *transaction) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	return getDB(ctx).Transaction(func(tx *gorm.DB) error {
		ctx = context.WithValue(ctx, config.DBKey, tx)

		return fn(ctx)
	})
}
