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

func FormCheckingRequest(existId int, creatorId int, formId int) error {

	queryChek := "SELECT id FROM forms WHERE id  = $1 and creator_id = $2"
	err := Db.QueryRow(queryChek, creatorId, formId).Scan(&existId)
	if err != nil {
		log.Printf("Ошибка при запросе на проверку наличия формы: %v", err)
		return err
	}
	return err
}

func FormCreateRequest(form models.FormRequest, creatorId int) (int, string, error) {

	link := uuid.New().String()
	query := `INSERT INTO forms (creator_id, title, description, link, private_key, expires_at) 
			  VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	var formId int
	err := Db.QueryRow(query, creatorId, form.Title, form.Description, link, form.PrivateKey, form.ExpiresAt).Scan(&formId)
	if err != nil {
		log.Printf("Ошибка при запросе создания формы: %v", err)
		return formId, link, err
	}
	return formId, link, err
}
func FormDeleteRequest(formId int, creatorId int) error {

	query := "DELETE FROM forms WHERE id = $1 AND creator_id = $2"
	_, err := Db.Exec(query, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе удаления формы: %v", err)
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
		log.Printf("Ошибка при запросе получения формы: %v", err)
		return form, err
	}
	return form, err
}

func FormUpdateRequest(updateForm models.FormRequest, creatorId int, formId int) error {

	query := "UPDATE forms SET title = $1, description = $2, private_key = $3, expires_at = $4 WHERE id = $5 AND creator_id = $6"
	_, err := Db.Exec(query, updateForm.Title, updateForm.Description, updateForm.PrivateKey, updateForm.ExpiresAt, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе обновления формы: %v", err)
		return err
	}
	return err
}

func GetFormsRequest(creatorId int) (*sql.Rows, error) {
	query := `SELECT id, title, description, link, private_key, expires_at FROM forms WHERE creator_id = $1 ORDER BY number_order`

	rows, err := Db.Query(query, creatorId)

	if err != nil {
		log.Printf("Не удалось найти формы через запрос: %v", err)
		return nil, err
	}

	return rows, err
}

func QuestionChekingRequest(existId int, creatorId int, formId int, questionId int) error {

	queryChek := "SELECT id FROM questions WHERE id = $1 AND form_id = $2 AND creator_id = $3"
	err := Db.QueryRow(queryChek, questionId, formId, creatorId).Scan(&existId)
	if err != nil {
		log.Printf("Ошибка при запросе на проверку наличия вопроса: %v", err)
		return err
	}
	return err
}
func QuestionCreateRequest(question models.QuestionRequest, formId int) (int, error) {
	query := `INSERT INTO questions (form_id, title, number_order, required) 
			  VALUES($1, $2, $3, $4) RETURNING id`

	var questionId int
	err := Db.QueryRow(query, formId, question.Title, question.NumberOrder, question.Required).Scan(&questionId)
	if err != nil {
		log.Printf("Ошибка при запросе создания вопроса: %v", err)
		return questionId, err
	}
	return questionId, err
}
func QuestionDeleteRequest(creator_id int, formId int, questionId int) (sql.Result, error) {

	query := "DELETE FROM questions WHERE id = $1 AND form_id = $2 and creator_id = $3"
	_, err := Db.Exec(query, questionId, formId)
	if err != nil {
		log.Printf("Ошибка при запросе удаления вопроса: %v", err)
		return nil, err
	}
	return nil, err
}

func QuestionsGetRequest(creator_id int, formId int) (*sql.Rows, error) {
	query := `SELECT questions.id, questions.title, questions.number_order, questions.required, answers.title, answers.number_order, answers.count 
			  FROM questions
			  JOIN answers ON questions.id = answers.question_id 
			  WHERE form_id = $1 AND creator_id = $2 ORDER BY questions.number_order, answers.number_order`

	rows, err := Db.Query(query, formId, creator_id)

	if err != nil {
		log.Printf("Ошибка при запросе получения вопросов: %v", err)
		return nil, err
	}
	return rows, err
}

func QuestionUpdateRequest(updateQuestion models.QuestionRequest, creator_id int, formId int, questionId int) error {

	query := "UPDATE questions SET title = $1, number_order = $2, required = $3 WHERE id = $4 AND form_id = $5 and creator_id = $6"
	_, err := Db.Exec(query, updateQuestion.Title, updateQuestion.NumberOrder, updateQuestion.Required, questionId, formId, creator_id)
	if err != nil {
		log.Printf("Ошибка при запросе обновления вопроса: %v", err)
		return err
	}
	return err
}

func AnswerChekingRequest(existId int, creatorId int, formId int, questionId int, answerId int) error {
	queryChek := "SELECT id FROM answers WHERE id = $1 AND question_id = $2 AND form_id = $3 AND creator_id = $4"
	err := Db.QueryRow(queryChek, answerId, questionId, formId, creatorId).Scan(&existId)
	if err != nil {
		log.Printf("Ошибка при запросе на проверку наличия ответа: %v", err)
		return err
	}
	return err
}

func AnswerCreateRequest(answer models.AnswerRequest, questionId int) (int, error) {
	query := `INSERT INTO answers (question_id, title, number_order, count) 
			  VALUES($1, $2, $3, $4) RETURNING id`

	var answerId int
	err := Db.QueryRow(query, questionId, answer.Title, answer.NumberOrder, answer.Count).Scan(&answerId)
	if err != nil {
		log.Printf("Ошибка при запросе создания ответа: %v", err)
		return answerId, err
	}
	return answerId, err
}
func AnswerDeleteRequest(creator_id int, formId int, questionId int, answerId int) (sql.Result, error) {
	query := "DELETE FROM answers WHERE id = $1 AND question_id = $2 AND form_id = $3 and creator_id = $4"
	_, err := Db.Exec(query, answerId, questionId, formId, creator_id)
	if err != nil {
		log.Printf("Ошибка при запросе удаления ответа: %v", err)
		return nil, err
	}
	return nil, err
}

func AnswerUpdateRequest(updateAnswer models.AnswerRequest, creator_id int, formId int, questionId int, answerId int) error {
	query := "UPDATE answers SET title = $1, number_order = $2, count = $3 WHERE id = $4 AND question_id = $5 AND form_id = $6 and creator_id = $7"
	_, err := Db.Exec(query, updateAnswer.Title, updateAnswer.NumberOrder, updateAnswer.Count, answerId, questionId, formId, creator_id)
	if err != nil {
		log.Printf("Ошибка при запросе обновления ответа: %v", err)
		return err
	}
	return err
}

func GetAnswersRequest(creator_id int, formId int, questionId int) (*sql.Rows, error) {
	query := `SELECT id, question_id, title, number_order, count 
			  FROM answers 
			  WHERE question_id = $1 AND form_id = $2 AND creator_id = $3`

	rows, err := Db.Query(query, questionId, formId, creator_id)

	if err != nil {
		log.Printf("Ошибка при запросе получения ответов: %v", err)
		return nil, err
	}
	return rows, err
}
