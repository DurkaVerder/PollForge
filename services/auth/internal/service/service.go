package service

import (
	// "auth/internal/kafka"
	"auth/internal/kafka"
	"auth/internal/models"
	"auth/internal/storage"
	"errors"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

const (
	userRegisteredEvent = "user_registered"
	userLoginEvent      = "user_login"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func handleRegistration(req models.UserRequest) (string, error) {
	// Проверяем, есть ли уже пользователь
	if err := CheckUserRequest(req); err != nil {
		return "", err
	}
	// Регистрируем и получаем токен (Jwt + Kafka внутри)
	return registerUserInternal(req)
}

func handleLogin(req models.UserRequest) (string, error) {
	// Логиним и получаем токен
	return loginUserInternal(req)
}

func GenerateJwt(userId string) (string, error) {

	claims := jwt.MapClaims{
		"id":  userId,
		"exp": time.Now().Add(time.Hour * 6).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func CheckUserRequest(request models.UserRequest) error {
	var exist bool

	// Потому что только одного надо выбрать если есть аккаунт c такой же почтой
	err := storage.Db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)", request.Email).Scan(&exist)
	if err != nil {
		log.Printf("Ошибка запроса к базе данных: %v", err.Error())
		return err
	}
	if exist {
		log.Printf("Такой пользователь уже есть, %v", exist)
		return errors.New("Такой пользователь уже есть")
	}
	return nil
}

func registerUserInternal(request models.UserRequest) (string, error) {
	var userId string
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), 10)
	if err != nil {
		log.Printf("Ошибка при создании хеша пароля")
		return "", err
	}

	// Потому что отправляем все данные для регистрации и возвращаем id для создания jwt
	userId, err = storage.Registration(hashedPassword, request)
	if err != nil {
		log.Printf("Ошибка при вводе данных в базу")
		log.Printf("%s", err.Error())
		return "", err
	}

	token, err := GenerateJwt(userId)
	if err != nil {
		log.Printf("Ошибка при создании токена")
		log.Printf("%s", err.Error())
		return "", err
	}

	kafkaMsg := models.MessageKafka{
		EventType: userRegisteredEvent,
		UserID:    userId,
	}
	if err := kafka.SendMessage(kafkaMsg); err != nil {
		log.Printf("Не удалось отправить сообщение Kafka: %v", err)
	}
	
	return token, err
}

func loginUserInternal(request models.UserRequest) (string, error) {
	userId, err := storage.CheckingLoggingData(request)
	if err != nil {
		log.Printf("Ошибка при сопоставлении пароля и почты")
		log.Printf("%s", err.Error())
		return "", err
	}

	kafkaMsg := models.MessageKafka{
		EventType: userLoginEvent,
		UserID:    userId,
	}

	if err := kafka.SendMessage(kafkaMsg); err != nil {
		log.Printf("Не удалось отправить сообщение Kafka: %v", err)
	}

	token, err := GenerateJwt(userId)
	return token, err
}
