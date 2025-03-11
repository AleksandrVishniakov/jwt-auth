package tokenizer

import (
	"errors"
	"time"

	"github.com/AleksandrVishniakov/jwt-auth/internal/handlers"
	"github.com/golang-jwt/jwt"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Tokenizer struct {
	signature []byte
	tokenTTL  time.Duration
}

func New(signature []byte, tokenTTL time.Duration) *Tokenizer {
	return &Tokenizer{
		signature: signature,
		tokenTTL:  tokenTTL,
	}
}

type tokenClaims struct {
	jwt.StandardClaims
	handlers.TokenData
}

func (t *Tokenizer) Token(userID int32, role string, permissionMask int64) (string, error) {
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  time.Now().Unix(),
			ExpiresAt: time.Now().Add(t.tokenTTL).Unix(),
		},
		TokenData: handlers.TokenData{
			UserID:         userID,
			Role:           role,
			PermissionMask: permissionMask,
		},
	}).SignedString(t.signature)

	if err != nil {
		return "", err
	}

	return token, nil
}

func (t *Tokenizer) Parse(token string) (data handlers.TokenData, err error) {
	jwtToken, err := jwt.ParseWithClaims(
		token, &tokenClaims{},
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrInvalidToken
			}

			return t.signature, nil
		})

	if err != nil {
		return handlers.TokenData{}, err
	}

	claims, ok := jwtToken.Claims.(*tokenClaims)
	if !ok {
		return handlers.TokenData{}, ErrInvalidToken
	}

	return claims.TokenData, nil
}
