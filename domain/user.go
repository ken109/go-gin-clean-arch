package domain

import (
	"context"
	"github.com/rs/xid"
	"go-gin-clean-arch/domain/vobj"
	"go-gin-clean-arch/packages/util"
	"go-gin-clean-arch/resource/request"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) (xid.ID, error)
	GetByID(ctx context.Context, id xid.ID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByRecoveryToken(ctx context.Context, recoveryToken string) (*User, error)
	Update(ctx context.Context, user *User) error

	EmailExists(ctx context.Context, email string) (bool, error)
}

type User struct {
	SoftDeleteModel
	Email    string        `json:"email" validate:"required,email" gorm:"index;unique"`
	Password vobj.Password `json:"-"`

	RecoveryToken *vobj.RecoveryToken `json:"-" gorm:"index"`
}

func NewUser(ctx context.Context, dto *request.UserCreate) (*User, error) {
	user := User{
		Email:         dto.Email,
		RecoveryToken: vobj.NewRecoveryToken(""),
	}

	util.Validate(ctx, user)

	password, err := vobj.NewPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil {
		return nil, err
	}

	user.Password = *password

	return &user, nil
}

func (u *User) ResetPassword(ctx context.Context, dto *request.UserResetPassword) error {
	if !u.RecoveryToken.IsValid() {
		util.InvalidField(ctx, "RecoveryToken", "リカバリートークンが無効です")
		return nil
	}

	password, err := vobj.NewPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil {
		return err
	}

	u.Password = *password

	u.RecoveryToken.Clear()
	return nil
}
