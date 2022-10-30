package authservice

import (
	"doko/gvn-ultimate-bot/models"
	"time"

	"gopkg.in/dgrijalva/jwt-go.v3"
)

type AuthService interface {
	IssueToken(u models.User) (string, error)
	ParseToken(token string) (*Claims, error)
}

type authService struct {
	jwtSecret string
}

func NewAuthService(jwtSecret string) AuthService {
	return &authService{
		jwtSecret: jwtSecret,
	}
}

type Claims struct {
	Email string `json:"email"`
	ID    uint   `json:"id"`
	jwt.StandardClaims
}

func (auth *authService) IssueToken(u models.User) (string, error) {
	nowTime := time.Now()
	expireTime := nowTime.Add(24 * time.Hour) // 24 hours

	claims := Claims{
		u.Email,
		u.ID,
		jwt.StandardClaims{
			ExpiresAt: expireTime.Unix(),
			Issuer:    "GVN Ultimate Bot",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return tokenClaims.SignedString([]byte(auth.jwtSecret))
}

func (auth *authService) ParseToken(token string) (*Claims, error) {
	tokenClaims, err := jwt.ParseWithClaims(
		token,
		&Claims{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(auth.jwtSecret), nil
		},
	)

	if tokenClaims != nil {
		claims, ok := tokenClaims.Claims.(*Claims)
		if ok && tokenClaims.Valid {
			return claims, nil
		}
	}

	return nil, err
}
