package service

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
)

func GetParamFromJWT(tokenBearer, param string) (string, error) {
	token := strings.TrimPrefix(tokenBearer, "Bearer ")

	parsedToken, err := parsedToken(token)
	if err != nil {
		return "", err
	}

	valid, err := validToken(parsedToken)
	if err != nil || !valid {
		return "", err
	}

	claims, err := getClaims(parsedToken)
	if err != nil {
		return "", err
	}

	value, ok := claims[param].(string)
	if !ok {
		return "", fmt.Errorf("param %s not found in claims", param)
	}

	return value, nil
}

func parsedToken(tokenStr string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}

func validToken(token *jwt.Token) (bool, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return false, errors.New("token is expired")
		}
		return true, nil
	}
	return false, errors.New("invalid token claims")
}

func getClaims(token *jwt.Token) (jwt.MapClaims, error) {
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}
	return nil, errors.New("invalid token claims")
}
