package jwtgen

import (
	"github.com/golang-jwt/jwt/v5"
)

type JwtGenerator struct {
	secret []byte
}

func NewJwtGenerator(secret []byte) JwtGenerator {
	return JwtGenerator{
		secret: secret,
	}
}

func (jwtgen JwtGenerator) GenToken(values jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, values)

	tokenString, err := token.SignedString(jwtgen.secret)
	return tokenString, err
}

func (jwtSecret JwtGenerator) VerifyToken(tokenstring string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenstring, func(t *jwt.Token) (interface{}, error) {
		return jwtSecret.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return token.Claims.(jwt.MapClaims), nil
}
