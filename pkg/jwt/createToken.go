package jwt

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaim struct {
	Type string `json:"type"`
	jwt.RegisteredClaims
}

var JWT_SECRET = []byte(os.Getenv(("JWT_SECRET")))

func CreateToken(user_id, token_type string, duration uint) (string, error) {

	claims := CustomClaim{
		Type: token_type,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user_id,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(duration))),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.NewString(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(JWT_SECRET)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
