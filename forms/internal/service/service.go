package service

import (
	"database/sql"
	"fmt"
	"forms/internal/models"
	"forms/internal/storage"
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt"
)

func GetToken(auth string) (*jwt.Token, error) {
	tokenStr := strings.TrimPrefix(auth, "Bearer")
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("Неподходящий метод подписи")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	return token, err
}

func FormChek(creatorId int,formId int)error{
	var existId int
	err := storage.FormChekingRequest(existId, creatorId, formId)
	if err != nil{
		log.Print("Ошибка при проверке на наличие формы")
		log.Printf("%s", err.Error())
		return err
	}
	return err
}

func FormDelete(formId int, creatorId int)(sql.Result,error){
	err := storage.FormDeleteRequest(formId, creatorId)
	if err != nil{
		log.Printf("Ошибка при удалении данных")
		return nil,err
	}
	return nil,err
}

func FormGet(creatorId int, formId int)(models.Form, error){
	form, err := storage.FormGetRequest(creatorId, formId)
	if err != nil{
		log.Printf("Ошибка при получении данных формы")
		return form, err
	}
	return form, err
}
func FormUpdate(updateForm models.FormRequest, creatorId int, formId int)(error){
	err := storage.FormUpdateRequest(updateForm, creatorId, formId)
	if err != nil{
		log.Printf("Ошибка при удалении данных")
		return err
	}
	return err
}

func FormCreate(form models.FormRequest, creatorId int)(int, string, error){
	formId, link, err := storage.FormCreateRequest(form, creatorId)
	if err != nil{
		log.Printf("Ошибка при удалении данных")
		return formId, link, err
	}
	return formId, link, err
}