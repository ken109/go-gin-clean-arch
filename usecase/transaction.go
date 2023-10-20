package usecase

import (
	"context"
)

type TransactionRepository interface {
	Do(ctx context.Context, fn func(ctx context.Context) error) error
}
