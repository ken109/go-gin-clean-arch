package repository

import (
	"go-gin-ddd/packages/context"

	"go-gin-ddd/domain/entity"
)

type IUser interface {
	Create(ctx context.Context, user *entity.User) (uint, error)
	GetByID(ctx context.Context, id uint) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	GetByRecoveryToken(ctx context.Context, recoveryToken string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error

	EmailExists(ctx context.Context, email string) (bool, error)
}
