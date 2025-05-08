package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Нет токена"})
			return
		}
		token, err := getToken(auth)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Токен не валиден"})
			return
		}
		if valid, err := validToken(token); !valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		claims := token.Claims.(jwt.MapClaims)
		rawId, ok := claims["id"]
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Нет id в токене"})
			return
		}
		c.Set("id", rawId)
		c.Next()
	}
}

func getToken(auth string) (*jwt.Token, error) {
	tokenStr := strings.TrimSpace(strings.TrimPrefix(auth, "Bearer"))
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("неподходящий метод подписи")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return token, err
}

func validToken(token *jwt.Token) (bool, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return false, errors.New("токен истёк")
		}
		return true, nil
	}
	return false, errors.New("токен не валиден")
}
