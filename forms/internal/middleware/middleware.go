package middleware

import (
	"fmt"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer") {
			c.AbortWithStatusJSON(401, gin.H{"error": "Нет токена"})
			return
		}

		tokenStr := strings.TrimPrefix(auth, "Bearer")
		token, err := jwt.Parse(tokenStr,func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256{
				return nil, fmt.Errorf("Неподходящий метод подписи")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil{
			c.AbortWithStatusJSON(401,gin.H{"error":"Токен не валиден"})
		}
		claims := token.Claims.(jwt.MapClaims)
		rawId, ok := claims["id"]
		if !ok{
			c.AbortWithStatusJSON(401,gin.H{"error":"Нет id в токене"})
		}
		c.Set("id", rawId)
		c.Next()
	}
}
