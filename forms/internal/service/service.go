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

func FormChek(creatorId int, formId int) error {
	var existId int
	err := storage.FormChekingRequest(existId, creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие формы: %v", err)
		log.Printf("%s", err.Error())
		return err
	}
	return err
}

func FormDelete(formId int, creatorId int) (sql.Result, error) {
	err := storage.FormDeleteRequest(formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return nil, err
	}
	return nil, err
}

func FormGet(creatorId int, formId int) (models.Form, error) {
	form, err := storage.FormGetRequest(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при получении данных формы: %v", err)
		return form, err
	}
	return form, err
}
func FormUpdate(updateForm models.FormRequest, creatorId int, formId int) error {
	err := storage.FormUpdateRequest(updateForm, creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return err
	}
	return err
}

func FormCreate(form models.FormRequest, creatorId int) (int, string, error) {
	formId, link, err := storage.FormCreateRequest(form, creatorId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return formId, link, err
	}
	return formId, link, err
}

func FormsGet(creatorId int) ([]models.Form, error) {
	rows, err := storage.GetFormsRequest(creatorId)
	var forms []models.Form
	if err != nil {
		log.Printf("Ошибка при получении форм: %v", err)
		return forms, err
	}
	defer rows.Close()
	for rows.Next() {
		var form models.Form
		err := rows.Scan(&form.Id,
			&form.Title,
			&form.Description,
			&form.Link,
			&form.PrivateKey,
			&form.ExpiresAt)
		if err != nil {
			log.Printf("Не удалось считать данные формы через запрос: %v", err)
			return forms, err
		}
		forms = append(forms, form)
	}

	return forms, err
}

func QuestionChek(creator_id int, formId int, questionId int) error {
	var existId int
	err := storage.QuestionChekingRequest(existId, creator_id, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие вопроса: %v", err)
		log.Printf("%s", err.Error())
		return err
	}
	return err
}

func QuestionDelete(creator_id int, formId int, questionId int) (sql.Result, error) {
	_, err := storage.QuestionDeleteRequest(creator_id, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return nil, err
	}
	return nil, err
}

func QuestionGet(creator_id int, formId int, questionId int) (models.Question, error) {
	question, err := storage.QuestionGetRequest(creator_id, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при получении данных вопроса: %v", err)
		return question, err
	}
	return question, err
}

func QuestionUpdate(updateQuestion models.QuestionRequest, creator_Id int, formId int, questionId int) error {
	err := storage.QuestionUpdateRequest(updateQuestion, creator_Id, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return err
	}
	return err
}
func QuestionCreate(question models.QuestionRequest, formId int) (int, error) {
	questionId, err := storage.QuestionCreateRequest(question, formId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return questionId, err
	}
	return questionId, err
}
func QuestionsGet(creator_Id, formId int) ([]models.Question, error) {
	rows, err := storage.QuestionsGetRequest(creator_Id, formId)
	var questions []models.Question
	if err != nil {
		log.Printf("Ошибка при получении вопросов: %v", err)
		return questions, err
	}
	defer rows.Close()
	for rows.Next() {
		var question models.Question
		err := rows.Scan(&question.Id,
			&question.FormId,
			&question.NumberOrder,
			&question.Title,
			&question.Required)
		if err != nil {
			log.Printf("Не удалось считать данные вопроса через запрос: %v", err)
			return questions, err
		}
		questions = append(questions, question)
	}

	return questions, err
}

func AnswerChek(creator_id int, formId int, questionId int, answerId int) error {
	var existId int
	err := storage.AnswerChekingRequest(existId, creator_id, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие ответа: %v", err)
		log.Printf("%s", err.Error())
		return err
	}
	return err
}

func AnswerDelete(creator_Id int, formId int, questionId int, answerId int) (sql.Result, error) {
	_, err := storage.AnswerDeleteRequest(creator_Id, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return nil, err
	}
	return nil, err
}

func AnswerGet(creator_Id int, formId int, questionId int, answerId int) (models.Answer, error) {
	answer, err := storage.AnswerGetRequest(creator_Id, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при получении данных ответа: %v", err)
		return answer, err
	}
	return answer, err
}

func AnswerUpdate(updateAnswer models.AnswerRequest, creator_Id int, formId int, questionId int, answerId int) error {
	err := storage.AnswerUpdateRequest(updateAnswer, creator_Id, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return err
	}
	return err
}

func AnswerCreate(answer models.AnswerRequest, formId int, questionId int) (int, error) {
	answerId, err := storage.AnswerCreateRequest(answer, questionId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return answerId, err
	}
	return answerId, err
}

func AnswersGet(creator_Id, formId int, questionId int) ([]models.Answer, error) {
	rows, err := storage.GetAnswersRequest(creator_Id, formId, questionId)
	var answers []models.Answer
	if err != nil {
		log.Printf("Ошибка при получении ответов: %v", err)
		return answers, err
	}
	defer rows.Close()
	for rows.Next() {
		var answer models.Answer
		err := rows.Scan(&answer.Id,
			&answer.QuestionId,
			&answer.Title,
			&answer.NumberOrder,
			&answer.Count,
			&answer.AnswerId)
		if err != nil {
			log.Printf("Не удалось считать данные ответа через запрос: %v", err)
			return answers, err
		}
		answers = append(answers, answer)
	}

	return answers, err
}
