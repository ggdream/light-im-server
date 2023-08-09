package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenMalformed    = jwt.ErrTokenMalformed
	ErrTokenExpired      = jwt.ErrTokenExpired
	ErrTokenNotValidYet  = jwt.ErrTokenNotValidYet
	ErrTokenSignedMethod = errors.New("token is not signed with hmac")
	// ErrTokenTypeNotValid = errors.New("field type of token is not valid")
)

type customClaims struct {
	UserID string `json:"uid"`
	jwt.RegisteredClaims
}

func Generate(key []byte, userId string) (string, int64, error) {
	claims := customClaims{
		userId,
		jwt.RegisteredClaims{
			Issuer:    "light-im",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	tk := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tkStr, err := tk.SignedString(key)
	if err != nil {
		return "", 0, err
	}

	return tkStr, claims.ExpiresAt.UnixMilli(), nil
}

func Parse(key []byte, tkStr string) (string, error) {
	tk, err := jwt.Parse(tkStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrTokenSignedMethod
		}
		return key, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := tk.Claims.(jwt.MapClaims)
	if ok && tk.Valid {
		// 校验成功
		return claims["uid"].(string), nil
	}

	return "", tk.Claims.Valid()
}
