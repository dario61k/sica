package jwt

import (
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

type MyClaims struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

func ValidateToken(s_token string) (*MyClaims, error) {

	token, err := jwt.ParseWithClaims(s_token, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, err
	}

	claims, ok := token.Claims.(*MyClaims)

	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
