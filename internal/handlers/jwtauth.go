package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/AleksandrVishniakov/jwt-auth/internal/e"
)

type TokenData struct {
	UserID         int32  `json:"userID"`
	Role           string `json:"role"`
	PermissionMask int64  `json:"permissionMask"`
}

type TokenParser interface {
	Parse(token string) (TokenData, error)
}

func JWTAuth(
	log *slog.Logger,
	parser TokenParser,
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := getTokenFromAuthHeader(r.Header.Get("Authorization"))
			if err != nil {
				httpError := e.Authorization(e.WithError(err))
				if err = EncodeResponse(w, httpError, httpError.Code); err != nil {
					slog.Error("encoding response error", e.SlogErr(err))
				}
				return
			}

			data, err := parser.Parse(token)
			if err != nil {
				httpError := e.Authorization(e.WithError(err))
				if err = EncodeResponse(w, httpError, httpError.Code); err != nil {
					slog.Error("encoding response error", e.SlogErr(err))
				}
				return
			}

			ctx = context.WithValue(ctx, userIDKey, data.UserID)
			ctx = context.WithValue(ctx, roleKey, data.Role)
			ctx = context.WithValue(ctx, permissionMaskKey, data.PermissionMask)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func getTokenFromAuthHeader(header string) (token string, err error) {
	const bearerAuthType = "Bearer"

	parts := strings.Split(header, " ")
	if len(parts) != 2 {
		return "", errors.New("invalid auth header length")
	}

	if parts[0] != bearerAuthType {
		return "", errors.New("invalid auth type")
	}

	if parts[1] == "" {
		return "", errors.New("empty token")
	}

	return parts[1], nil
}
