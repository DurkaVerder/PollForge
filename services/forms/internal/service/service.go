package service

import (
	"database/sql"
	"forms/internal/models"
	"forms/internal/storage"
	"log"
)

func FormCheck(creatorId int, formId int) error {
	var existId int
	err := storage.FormCheckingRequest(existId, creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие формы: %v", err)
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

func FormGet(creatorId, formId int) (models.Form, error) {
	var form models.Form

	form, err := storage.FormGetRequest(creatorId, formId)
	if err != nil {
		return form, err
	}
	form.CreatorId = creatorId

	// Получение вопросов с ответами
	form.Questions, err = storage.QuestionsWithAnswersGet(formId, creatorId)
	if err != nil {
		return form, err
	}

	return form, nil
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


func QuestionChek(creator_id int, formId int, questionId int) error {
	var existId int
	err := storage.QuestionChekingRequest(existId, creator_id, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие вопроса: %v", err)
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

func QuestionUpdate(updateQuestion models.QuestionRequest, creator_Id int, formId int, questionId int) error {
	err := storage.QuestionUpdateRequest(updateQuestion, creator_Id, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при изменении данных: %v", err)
		return err
	}
	return err
}
func QuestionCreate(question models.QuestionRequest, creatorId int, formId int) (int, error) {
	questionId, err := storage.QuestionCreateRequest(question, creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при создании вопроса: %v", err)
		return questionId, err
	}
	return questionId, err
}

func AnswerChek(creator_id int, formId int, questionId int, answerId int) error {
	var existId int
	err := storage.AnswerChekingRequest(existId, creator_id, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при проверке на наличие ответа: %v", err)
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

func AnswerUpdate(updateAnswer models.AnswerRequest, creator_Id int, formId int, questionId int, answerId int) error {
	err := storage.AnswerUpdateRequest(updateAnswer, creator_Id, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return err
	}
	return err
}

func AnswerCreate(answer models.AnswerRequest, creatorId int, formId int, questionId int) (int, error) {
	answerId, err := storage.AnswerCreateRequest(answer, creatorId, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при удалении данных: %v", err)
		return answerId, err
	}
	return answerId, err
}

func GetFormByLink(link string) (models.Form, error) {
	var form models.Form

	form, err := storage.GetFormByLinkRequest(link)
	if err != nil {
		log.Printf("Ошибка при получении формы по ссылке: %v", err)
		return form, err
	}

	// Получение вопросов с ответами
	form.Questions, err = storage.QuestionsWithAnswersGet(form.Id, form.CreatorId)
	if err != nil {
		return form, err
	}

	return form, nil
}