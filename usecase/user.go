package usecase

import (
	"net/http"

	"context"
	jwt "github.com/ken109/gin-jwt"
	"github.com/rs/xid"
	"go-gin-clean-arch/adapter/gateway/mail"
	"go-gin-clean-arch/config"
	"go-gin-clean-arch/domain"
	"go-gin-clean-arch/packages/errors"
	"go-gin-clean-arch/packages/util"
	"go-gin-clean-arch/resource/mailbody"
	"go-gin-clean-arch/resource/request"
	"go-gin-clean-arch/resource/response"
)

type User interface {
	Create(ctx context.Context, req *request.UserCreate) (xid.ID, error)

	ResetPasswordRequest(ctx context.Context, req *request.UserResetPasswordRequest) (*response.UserResetPasswordRequest, error)
	ResetPassword(ctx context.Context, req *request.UserResetPassword) error
	Login(ctx context.Context, req *request.UserLogin) (*response.UserLogin, error)
	RefreshToken(req *request.UserRefreshToken) (*response.UserLogin, error)

	GetByID(ctx context.Context, id xid.ID) (*domain.User, error)
}

type user struct {
	transactionRepo domain.TransactionRepository
	userRepo        domain.UserRepository
	email           mail.Sender
}

func NewUser(txr domain.TransactionRepository, tr domain.UserRepository, email mail.Sender) User {
	return &user{
		transactionRepo: txr,
		userRepo:        tr,
		email:           email,
	}
}

func (u user) Create(ctx context.Context, req *request.UserCreate) (xid.ID, error) {
	emailExists, err := u.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return xid.NilID(), err
	}

	if emailExists {
		util.InvalidField(ctx, "email", "既に使用されています")
	}

	newUser, err := domain.NewUser(ctx, req)
	if err != nil {
		return xid.NilID(), err
	}

	if util.IsInvalid(ctx) {
		return xid.NilID(), util.ValidationError(ctx)
	}

	id, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return xid.NilID(), err
	}

	return id, nil
}

func (u user) ResetPasswordRequest(ctx context.Context, req *request.UserResetPasswordRequest) (*response.UserResetPasswordRequest, error) {
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

	err = u.transactionRepo.Do(ctx, func(ctx context.Context) error {
		err = u.userRepo.Update(ctx, user)
		if err != nil {
			return err
		}

		err = u.email.Send(user.Email, mailbody.UserResetPasswordRequest{
			URL:   config.Env.App.URL,
			Token: user.RecoveryToken.String(),
		})
		if err != nil {
			return err
		}

		return nil
	})

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

	if util.IsInvalid(ctx) {
		return util.ValidationError(ctx)
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

func (u user) RefreshToken(req *request.UserRefreshToken) (*response.UserLogin, error) {
	var (
		res response.UserLogin
		ok  bool
		err error
	)

	ok, res.Token, res.RefreshToken, err = jwt.RefreshToken(config.UserRealm, req.RefreshToken)
	if err != nil {
		return nil, errors.NewUnexpected(err)
	}

	if !ok {
		return nil, nil
	}
	return &res, nil
}

func (u user) GetByID(ctx context.Context, id xid.ID) (*domain.User, error) {
	res, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res, nil
}
