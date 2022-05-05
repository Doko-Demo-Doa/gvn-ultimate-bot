package middlewares

import (
	"strings"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

type Claims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

func stripBearer(tok string) (string, error) {
	if len(tok) > 6 && strings.ToLower(tok[0:7]) == "bearer" {
		return tok[7:], nil

	}

	return tok, nil
}
