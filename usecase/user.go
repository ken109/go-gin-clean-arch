package usecase

import (
	"net/http"

	jwt "github.com/ken109/gin-jwt"
	"github.com/rs/xid"
	"go-gin-clean-arch/domain"
	"go-gin-clean-arch/resource/mail_body"

	"go-gin-clean-arch/packages/context"
	"go-gin-clean-arch/packages/errors"

	"go-gin-clean-arch/config"
	"go-gin-clean-arch/resource/request"
	"go-gin-clean-arch/resource/response"
)

type UserInputPort interface {
	Create(ctx context.Context, req *request.UserCreate) error

	ResetPasswordRequest(ctx context.Context, req *request.UserResetPasswordRequest) error
	ResetPassword(ctx context.Context, req *request.UserResetPassword) error
	Login(ctx context.Context, req *request.UserLogin) error
	RefreshToken(req *request.UserRefreshToken) error

	GetByID(ctx context.Context, id xid.ID) error
}

type UserOutputPort interface {
	Create(id xid.ID) error

	ResetPasswordRequest(res *response.UserResetPasswordRequest) error
	ResetPassword() error
	Login(isSession bool, res *response.UserLogin) error
	RefreshToken(isSession bool, res *response.UserLogin) error

	GetByID(res *domain.User) error
}

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) (xid.ID, error)
	GetByID(ctx context.Context, id xid.ID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	GetByRecoveryToken(ctx context.Context, recoveryToken string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error

	EmailExists(ctx context.Context, email string) (bool, error)
}

type user struct {
	outputPort UserOutputPort
	userRepo   UserRepository
	email      Mail
}

type UserInputFactory func(outputPort UserOutputPort) UserInputPort

func NewUserInputFactory(tr UserRepository, email Mail) UserInputFactory {
	return func(o UserOutputPort) UserInputPort {
		return &user{
			outputPort: o,
			userRepo:   tr,
			email:      email,
		}
	}
}

func (u user) Create(ctx context.Context, req *request.UserCreate) error {
	email, err := u.userRepo.EmailExists(ctx, req.Email)
	if err != nil {
		return err
	}

	if email {
		ctx.FieldError("Email", "既に使用されています")
	}

	newUser, err := domain.NewUser(ctx, req)
	if err != nil {
		return err
	}

	if ctx.IsInValid() {
		return ctx.ValidationError()
	}

	id, err := u.userRepo.Create(ctx, newUser)
	if err != nil {
		return err
	}

	return u.outputPort.Create(id)
}

func (u user) ResetPasswordRequest(ctx context.Context, req *request.UserResetPasswordRequest) error {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		switch v := err.(type) {
		case *errors.Expected:
			if !v.ChangeStatus(http.StatusNotFound, http.StatusOK) {
				return err
			}
		default:
			return err
		}
	}

	var res response.UserResetPasswordRequest

	res.Duration, res.Expire, err = user.RecoveryToken.Generate()
	if err != nil {
		return err
	}

	err = ctx.Transaction(
		func(ctx context.Context) error {
			err = u.userRepo.Update(ctx, user)
			if err != nil {
				return err
			}

			err = u.email.Send(user.Email, mail_body.UserResetPasswordRequest{
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
		return err
	}

	return u.outputPort.ResetPasswordRequest(&res)
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

func (u user) Login(ctx context.Context, req *request.UserLogin) error {
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return err
	}

	if user.Password.IsValid(req.Password) {
		var res response.UserLogin

		res.Token, res.RefreshToken, err = jwt.IssueToken(config.UserRealm, jwt.Claims{
			"uid": user.ID,
		})
		if err != nil {
			return errors.NewUnexpected(err)
		}
		return u.outputPort.Login(req.Session, &res)
	}
	return u.outputPort.Login(req.Session, nil)
}

func (u user) RefreshToken(req *request.UserRefreshToken) error {
	var (
		res response.UserLogin
		ok  bool
		err error
	)

	ok, res.Token, res.RefreshToken, err = jwt.RefreshToken(config.UserRealm, req.RefreshToken)
	if err != nil {
		return errors.NewUnexpected(err)
	}

	if !ok {
		return nil
	}
	return u.outputPort.RefreshToken(req.Session, &res)
}

func (u user) GetByID(ctx context.Context, id xid.ID) error {
	res, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return u.outputPort.GetByID(res)
}
