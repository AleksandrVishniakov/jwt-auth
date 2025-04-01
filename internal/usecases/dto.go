package usecases

import (
	"errors"
	"slices"
	"unicode/utf8"
)

type LoginRequest struct {
	Login    string
	Password string
}

type LoginResponse struct {
	ID    int32
	Token string
}

func NewLoginRequest(
	login string,
	password string,
) (*LoginRequest, error) {
	if len(login) < 3 || len(login) > 64 {
		return nil, errors.New("invalid login length")
	}

	if !utf8.ValidString(login) {
		return nil, errors.New("login is invalid string")
	}

	if len(login) < 3 || len(password) > 128 {
		return nil, errors.New("invalid password length")
	}

	return &LoginRequest{
		Login:    login,
		Password: password,
	}, nil
}

type RegisterRequest struct {
	Login    string
	Password string
}

type RegisterResponse struct {
	ID    int32
	Token string
}

func NewRegisterRequest(
	login string,
	password string,
) (*RegisterRequest, error) {
	if len(login) < 3 || len(login) > 64 {
		return nil, errors.New("invalid login length")
	}

	if !utf8.ValidString(login) {
		return nil, errors.New("login is invalid string")
	}

	if len(login) < 3 || len(password) > 128 {
		return nil, errors.New("invalid password length")
	}

	return &RegisterRequest{
		Login:    login,
		Password: password,
	}, nil
}

type UpdateUserRoleRequest struct {
	UserID          int32
	Role            string
	PermissionsMask int64
}

func NewUpdateUserRoleRequest(
	userID int32,
	role string,
	permissionsMask int64,
) (*UpdateUserRoleRequest, error) {
	var existingRoles = []string{"student", "admin"}

	if userID < 1 {
		return nil, errors.New("invalid user id")
	}

	if permissionsMask < 0 {
		return nil, errors.New("invalid permissions mask")
	}

	if !slices.Contains(existingRoles, role) {
		return nil, errors.New("unknown role")
	}

	return &UpdateUserRoleRequest{
		UserID:          userID,
		Role:            role,
		PermissionsMask: permissionsMask,
	}, nil
}

type GetUserByIDRequest struct {
	UserID         int32
	PermissionMask int64
	ProfileID      int32
}

type GetUserByIDResponse struct {
	ID    int32
	Login string
	Role  string
}

func NewGetUserByIDRequest(
	userID int32,
	permissionMask int64,
	profileID int32,
) (*GetUserByIDRequest, error) {
	return &GetUserByIDRequest{
		UserID:         userID,
		PermissionMask: permissionMask,
		ProfileID:      profileID,
	}, nil
}
