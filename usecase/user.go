package usecase

import (
	"net/http"

	jwt "github.com/ken109/gin-jwt"

	"packages/context"
	"packages/errors"

	"go-gin-ddd/config"
	"go-gin-ddd/domain/entity"
	"go-gin-ddd/domain/repository"
	emailInfra "go-gin-ddd/infrastructure/email"
	"go-gin-ddd/resource/request"
	"go-gin-ddd/resource/response"
)

type IUser interface {
	Create(ctx context.Context, req *request.UserCreate) (uint, error)

	ResetPasswordRequest(
		ctx context.Context,
		req *request.UserResetPasswordRequest,
	) (*response.UserResetPasswordRequest, error)
	ResetPassword(ctx context.Context, req *request.UserResetPassword) error
	Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error)
	RefreshToken(refreshToken string) (*response.UserLogin, error)

	GetByID(ctx context.Context, id uint) (*entity.User, error)
}

type user struct {
	email    emailInfra.IEmail
	userRepo repository.IUser
}

func NewUser(email emailInfra.IEmail, tr repository.IUser) IUser {
	return &user{
		email:    email,
		userRepo: tr,
	}
}

func (u user) Create(ctx context.Context, req *request.UserCreate) (uint, error) {
	email, err := u.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return 0, err
	}

	if email {
		ctx.FieldError("Email", "既に使用されています")
	}

	newUser, err := entity.NewUser(ctx, req)
	if err != nil {
		return 0, err
	}

	if ctx.IsInValid() {
		return 0, ctx.ValidationError()
	}

	id, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (u user) ResetPasswordRequest(
	ctx context.Context,
	req *request.UserResetPasswordRequest,
) (*response.UserResetPasswordRequest, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		switch v := err.(type) {
		case *errors.Expected:
			if !v.ChangeStatus(http.StatusNotFound, http.StatusOK) {
				return nil, err
			}
		default:
			return nil, err
		}
	}

	var res response.UserResetPasswordRequest

	res.Duration, res.Expire, err = user.RecoveryToken.Generate()
	if err != nil {
		return nil, err
	}

	err = ctx.Transaction(
		func(ctx context.Context) error {
			err = u.userRepo.Update(ctx, user)
			if err != nil {
				return err
			}

			err = u.email.Send(user.Email, emailInfra.UserResetPasswordRequest{
				URL:   config.Env.App.URL,
				Token: user.RecoveryToken.String(),
			})
			if err != nil {
				return err
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (u user) ResetPassword(ctx context.Context, req *request.UserResetPassword) error {
	user, err := u.userRepo.GetByRecoveryToken(ctx, req.RecoveryToken)
	if err != nil {
		return err
	}

	err = user.ResetPassword(ctx, req)
	if err != nil {
		return err
	}

	if ctx.IsInValid() {
		return ctx.ValidationError()
	}

	return u.userRepo.Update(ctx, user)
}

func (u user) Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error) {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if user.Password.IsValid(req.Password) {
		var res response.UserLogin

		res.Token, res.RefreshToken, err = jwt.IssueToken(config.UserRealm, jwt.Claims{
			"uid": user.ID,
		})
		if err != nil {
			return nil, errors.NewUnexpected(err)
		}
		return &res, nil
	}
	return nil, nil
}

func (u user) RefreshToken(refreshToken string) (*response.UserLogin, error) {
	var (
		res response.UserLogin
		ok  bool
		err error
	)

	ok, res.Token, res.RefreshToken, err = jwt.RefreshToken(config.UserRealm, refreshToken)
	if err != nil {
		return nil, errors.NewUnexpected(err)
	}

	if !ok {
		return nil, nil
	}
	return &res, nil
}

func (u user) GetByID(ctx context.Context, id uint) (*entity.User, error) {
	return u.userRepo.GetByID(ctx, id)
}
