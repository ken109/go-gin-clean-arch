package mysql

import (
	"github.com/rs/xid"

	"context"
	"go-gin-clean-arch/domain"
	"go-gin-clean-arch/domain/vobj"
)

type user struct{}

func NewUser() domain.UserRepository {
	return &user{}
}

func (u user) Create(ctx context.Context, user *domain.User) (xid.ID, error) {
	db := getDB(ctx)

	if err := db.Create(user).Error; err != nil {
		return xid.NilID(), dbError(err)
	}

	return user.ID, nil
}

func (u user) GetByID(ctx context.Context, id xid.ID) (*domain.User, error) {
	db := getDB(ctx)

	var user domain.User
	err := db.First(&user, id).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &user, nil
}

func (u user) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	db := getDB(ctx)

	var dest domain.User
	err := db.Where(&domain.User{Email: email}).First(&dest).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &dest, nil
}

func (u user) GetByRecoveryToken(ctx context.Context, recoveryToken string) (*domain.User, error) {
	db := getDB(ctx)

	var dest domain.User
	err := db.Where(&domain.User{RecoveryToken: vobj.NewRecoveryToken(recoveryToken)}).First(&dest).Error
	if err != nil {
		return nil, dbError(err)
	}
	return &dest, nil
}

func (u user) Update(ctx context.Context, user *domain.User) error {
	db := getDB(ctx)

	return dbError(db.Model(user).Updates(user).Error)
}

func (u user) EmailExists(ctx context.Context, email string) (bool, error) {
	db := getDB(ctx)

	var count int64
	if err := db.Model(&domain.User{}).Where(map[string]any{"email": email}).Count(&count).Error; err != nil {
		return false, dbError(err)
	}
	return count > 0, nil
}
