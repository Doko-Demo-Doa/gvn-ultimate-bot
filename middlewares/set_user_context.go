package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/dgrijalva/jwt-go.v3"
)

// Set the "user" context if the user has valid token, otherwise do nothing
func SetUserContext(jwtSecret string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, _ := stripBearer(ctx.Request.Header.Get("Authorization"))

		tokenClaims, _ := jwt.ParseWithClaims(token, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(jwtSecret), nil
		})

		if tokenClaims != nil {
			claims, ok := tokenClaims.Claims.(*Claims)
			if ok && tokenClaims.Valid {
				// Set gin context values
			}
			fmt.Printf(claims.Subject)
		}
	}
}
