package storage

import (
	"database/sql"
	"forms/internal/models"
	"log"
	"os"

	"github.com/google/uuid"
)

var Db *sql.DB

func ConnectToDb() error {
	dsn := os.Getenv("DB_URL")
	var err error

	Db, err = sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	err = Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

func FormChekingRequest(existId int, creatorId int, formId int) error {

	queryChek := "SELECT id FROM forms WHERE id  = $1 and creator_id = $2"
	err := Db.QueryRow(queryChek, creatorId, formId).Scan(&existId)
	if err != nil {
		log.Printf("Ошибка при запросе на проверку наличия формы")
		return err
	}
	return err
}

func FormCreateRequest(form models.FormRequest, creatorId int) (int, string, error){

	link := uuid.New().String()
	query := `INSERT INTO forms (creator_id, title, description, link, private_key, expires_at) 
			  VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	var formId int
	err := Db.QueryRow(query, creatorId, form.Title, form.Description, link, form.PrivateKey, form.ExpiresAt).Scan(&formId)
	if err != nil{
		log.Printf("Ошибка при запросе создания формы")
		return formId, link, err
	}
	return formId, link, err
}
func FormDeleteRequest(formId int, creatorId int) error {

	query := "DELETE FROM forms WHERE id = $1 AND creator_id = $2"
	_, err := Db.Exec(query, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе удаления формы")
		return err
	}
	return err
}

func FormGetRequest(creatorId int, formId int) (models.Form, error) {

	query := `SELECT id, title, description, link, private_key, expires_at 
			  FROM forms 
			  WHERE id = $1 AND creator_id = $2`
	var form models.Form
	err := Db.QueryRow(query, formId, creatorId).Scan(&form.Id,
		&form.Title,
		&form.Description,
		&form.Link,
		&form.PrivateKey,
		&form.ExpiresAt)
	if err != nil {
		log.Printf("Ошибка при запросе получения формы")
		return form, err
	}
	return form, err
}

func FormUpdateRequest(updateForm models.FormRequest, creatorId int, formId int) error {

	query := "UPDATE forms SET title = $1, description = $2, private_key = $3, expires_at = $4 WHERE id = $5 AND creator_id = $6"
	_, err := Db.Exec(query, updateForm.Title, updateForm.Description, updateForm.PrivateKey, updateForm.ExpiresAt, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе обновления формы")
		return err
	}
	return err
}
