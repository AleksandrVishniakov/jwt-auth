package usecases

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/AleksandrVishniakov/jwt-auth/internal/e"
	"github.com/AleksandrVishniakov/jwt-auth/internal/roles"
	"golang.org/x/crypto/bcrypt"
)

type UserStorage interface {
	CreateUser(
		ctx context.Context,
		login string,
		passwordHash []byte,
	) (id int32, err error)

	CreateSuperUser(
		ctx context.Context,
		login string,
		passwordHash []byte,
	) (id int32, err error)

	GetUserById(
		ctx context.Context,
		id int32,
	) (user *UserModel, err error)

	GetUserByLogin(
		ctx context.Context,
		login string,
	) (user *UserModel, err error)

	UpdateRoleById(
		ctx context.Context,
		userID int32,
		roleAlias string,
	) (err error)
}

type TokenGenerator interface {
	Token(
		userID int32, 
		role string, 
		permissionMask int64,
	) (string, error)
}

type Usecase struct {
	log *slog.Logger
	storage UserStorage
	tokenGenerator TokenGenerator
}

func New(
	log *slog.Logger,
	storage UserStorage,
	tokenGenerator TokenGenerator,
) *Usecase {
	return &Usecase{
		log: log,
		storage: storage,
		tokenGenerator: tokenGenerator,
	}
}

func (u *Usecase) Login(
	ctx context.Context,
	req *LoginRequest,
) (resp *LoginResponse, err error) {
	const src = "Usecase.Login"
	log := u.log.With(slog.String("src", src))
	log.Debug("login user")

	user, err := u.storage.GetUserByLogin(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user: %w", src, err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, fmt.Errorf("%s: failed to compare password: %w", src, err)
	}

	token, err := u.tokenGenerator.Token(user.ID, user.Role, user.PermissionMask)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to generate token: %w", src, err)
	}

	return &LoginResponse{
		ID: user.ID,
		Token: token,
	}, nil
}

func (u *Usecase) Register(
	ctx context.Context,
	req *RegisterRequest,
) (resp *RegisterResponse, err error) {
	const src = "Usecase.Register"
	log := u.log.With(slog.String("src", src))
	log.Debug("register new user")

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to generate password hash: %w", src, err)
	}

	id, err := u.storage.CreateUser(ctx, req.Login, hash)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to create new user: %w", src, err)
	}

	user, err := u.storage.GetUserById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to get user role: %w", src, err)
	}

	token, err := u.tokenGenerator.Token(id, user.Role, user.PermissionMask)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to generate token: %w", src, err)
	}

	return &RegisterResponse{
		ID: id,
		Token: token,
	}, nil
}

func (u *Usecase) UpdateRoleById(
	ctx context.Context,
	req *UpdateUserRoleRequest,
) (err error) {
	const src = "Usecase.UpdateRoleById"
	log := u.log.With(slog.String("src", src))
	log.Debug("updating user role", slog.Int("id", int(req.UserID)))

	if !roles.HasPermission(req.PermissionsMask, roles.CanUpdateUserRole) {
		return e.ErrForbiddenAction
	}

	err = u.storage.UpdateRoleById(ctx, req.UserID, req.Role)
	if err != nil {
		return fmt.Errorf("%s: failed to update role: %w", src, err)
	}

	return nil
}

func (u *Usecase) CreateSuperUser(
	ctx context.Context,
	login string,
	password string,
) (err error) {
	const src = "Usecase.CreateSuperUser"
	log := u.log.With(slog.String("src", src))
	log.Debug("creating super user", slog.String("login", login))

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("%s: failed to generate password hash: %w", src, err)
	}

	_, err = u.storage.CreateSuperUser(ctx, login, hash)
	if err != nil && !errors.Is(err, e.ErrAlreadyExists) {
		return fmt.Errorf("%s: failed to create new user: %w", src, err)
	}

	return nil
}

