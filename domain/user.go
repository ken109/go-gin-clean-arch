package domain

import (
	"go-gin-ddd/packages/context"
	"go-gin-ddd/resource/response"

	"go-gin-ddd/domain/vobj"
	"go-gin-ddd/resource/request"
)

type User struct {
	softDeleteModel
	Email    string        `json:"email" validate:"required" gorm:"index;unique"`
	Password vobj.Password `json:"-"`

	RecoveryToken *vobj.RecoveryToken `json:"-" gorm:"index"`
}

type UserUsecase interface {
	Create(ctx context.Context, req *request.UserCreate) (uint, error)

	ResetPasswordRequest(
		ctx context.Context,
		req *request.UserResetPasswordRequest,
	) (*response.UserResetPasswordRequest, error)
	ResetPassword(ctx context.Context, req *request.UserResetPassword) error
	Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error)
	RefreshToken(refreshToken string) (*response.UserLogin, error)

	GetByID(ctx context.Context, id uint) (*User, error)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) (uint, error)
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByRecoveryToken(ctx context.Context, recoveryToken string) (*User, error)
	Update(ctx context.Context, user *User) error

	EmailExists(ctx context.Context, email string) (bool, error)
}

func NewUser(ctx context.Context, dto *request.UserCreate) (*User, error) {
	var user = User{
		Email:         dto.Email,
		RecoveryToken: vobj.NewRecoveryToken(""),
	}

	ctx.Validate(user)

	password, err := vobj.NewPassword(ctx, dto.Password, dto.PasswordConfirm)
	if err != nil {
		return nil, err
	}

	user.Password = *password

	return &user, nil
}

func (u *User) ResetPassword(ctx context.Context, dto *request.UserResetPassword) error {
	if !u.RecoveryToken.IsValid() {
		ctx.FieldError("RecoveryToken", "リカバリートークンが無効です")
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
