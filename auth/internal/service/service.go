package service

import (
	"auth/internal/models"
	"auth/internal/storage"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJwt(UserId string) (string, error) {

	claims := jwt.MapClaims{
		"sub": UserId,
		"exp": time.Now().Add(time.Hour * 6).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	return token.SignedString(jwtKey)
}

func CheckUserRequest(request models.RegisterRequest) error {
	var exist bool
	err := storage.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", request.Email).Scan(&exist) //Потому что только одного надо выбрать если есть похожий
	if err != nil {
		return err
	}
	if exist {
		return errors.New("Такой пользователь уже есть")
	}
	return nil
}

func RegistrUser(request models.RegisterRequest)(string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil{
		log.Printf("Ошибка при создании хеша пароля")
		return "", err
	}
	err = storage.InsertData(hashedPassword, request)
	if err != nil{
		log.Printf("Ошибка при вводе данных в базу данных")
	}	
}
