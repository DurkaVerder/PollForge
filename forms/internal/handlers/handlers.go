package handlers

import (
	"database/sql"
	"fmt"
	"forms/internal/models"
	"forms/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func extractUserID(c *gin.Context) (int, error) {
	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return 0, fmt.Errorf("id пользователя не найден")
	}
	creatorId, ok := creatorIdfl.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный тип id"})
		return 0, fmt.Errorf("неправильный тип id: %v", creatorIdfl)
	}
	return creatorId, nil
}

func extractFormID(c *gin.Context) (int, error) {
	formIdstr := c.Param("id")
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return 0, fmt.Errorf("неправильный тип id: %v", formIdstr)
	}
	return formId, nil
}

func extractQuestionID(c *gin.Context) (int, error) {
	questionIdstr := c.Param("question_id")
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return 0, fmt.Errorf("неправильный тип id: %v", questionIdstr)
	}
	return questionId, nil
}

func extractAnswerID(c *gin.Context) (int, error) {
	answerIdstr := c.Param("answer_id")
	answerId, err := strconv.Atoi(answerIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return 0, fmt.Errorf("неправильный тип id: %v", answerIdstr)
	}
	return answerId, nil
}

func CreateForm(c *gin.Context) {
	var form models.FormRequest

	err := c.ShouldBindJSON(&form)
	if err != nil {
		log.Printf("Ошибка при извлечении данных формы и переводе их в json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	var formId int
	formId, link, err := service.FormCreate(form, creatorId)

	if err != nil {
		log.Printf("Не удалось создать форму: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось создать форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"form_id": formId, "link": link})
}

func GetForm(c *gin.Context) {

	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := extractFormID(c)

	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}
	var form models.Form

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	form, err = service.FormGet(creatorId, formId)
	// Проверка на наличие возвращаемого значения, если форм нет - сработает условие цикла
	if err == sql.ErrNoRows {
		log.Printf("Не удалось получить формы из-за их отсутствия: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	} else if err != nil {
		log.Printf("Не удалось получить формы из-за ошибки: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка при получении формы"})
		return
	}

	c.JSON(http.StatusOK, form)
}

func GetForms(c *gin.Context) {

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	forms, err := service.FormsGet(creatorId)
	if err != nil {
		log.Printf("Не удалось получить формы: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось найти формы"})
		return
	}
	c.JSON(http.StatusOK, forms)
}

func UpdateForm(c *gin.Context) {

	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	var updateForm models.FormRequest

	err = c.ShouldBindJSON(&updateForm)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	err = service.FormCheck(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на существование формы: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	err = service.FormUpdate(updateForm, creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при обновлении данных формы: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка обновления данных"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно обновлена"})
}

func DeleteForm(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	// Проверка на существование формы для удаления, нужен id пользователя и id формы
	err = service.FormCheck(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на существование формы: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	// Удаление формы, для этого требуется id формы и пользователя из-за нужды нахождения формы для удаления с помощью sql-запроса
	_, err = service.FormDelete(formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при удалении формы: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно удалена"})
}

func CreateQuestion(c *gin.Context) {

	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	var question models.QuestionRequest

	err = c.ShouldBindJSON(&question)

	if err != nil {
		log.Printf("Ошибка при извлечении данных вопроса и переводе их в json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	err = service.FormCheck(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на существование формы: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	questionId, err := service.QuestionCreate(question, formId)
	if err != nil {
		log.Printf("Ошибка при создании вопросы: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось создать вопрос"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Вопрос успешно создан",
		"id вопроса": questionId})

}

func GetQuestions(c *gin.Context) {

	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	err = service.FormCheck(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при проверке на существование формы: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	questions, err := service.QuestionsGet(creatorId, formId)
	if err != nil {
		log.Printf("Ошибка при получении вопроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось получить вопросы"})
		return
	}
	c.JSON(http.StatusOK, questions)

}

func UpdateQuestion(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionId, err := extractQuestionID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	err = service.QuestionChek(creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Вопрос не найден"})
		return
	}

	var updateQuestion models.QuestionRequest

	err = c.ShouldBindJSON(&updateQuestion)

	if err != nil {
		log.Printf("Ошибка при извлечении данных вопроса для его обновления и переводе их в json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	err = service.QuestionUpdate(updateQuestion, creatorId, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при обновлении данных вопроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось обновить вопрос"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Вопрос успешно обновлен"})

}

func DeleteQuestion(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionId, err := extractQuestionID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	err = service.QuestionChek(creatorId, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при проверке вопросы на существование: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Вопрос не найден"})
		return
	}

	_, err = service.QuestionDelete(creatorId, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при удалении вопроса: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить вопрос"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Вопрос успешно удален"})

}

func CreateAnswer(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionId, err := extractQuestionID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	err = service.QuestionChek(creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Вопрос не найден"})
		return
	}

	var answer models.AnswerRequest

	err = c.ShouldBindJSON(&answer)

	if err != nil {
		log.Printf("Ошибка при извлечении данных ответа и переводе их в json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	answerId, err := service.AnswerCreate(answer, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при создании ответа: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось создать ответ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Ответ успешно создан",
		"id ответа": answerId})

}

func GetAnswers(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionId, err := extractQuestionID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	answerId, err := extractAnswerID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return
	}

	err = service.AnswerChek(creatorId, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при проверке ответа на существование: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Вопрос не найден"})
		return
	}

	var answers []models.Answer

	answers, err = service.AnswersGet(creatorId, formId, questionId)
	if err != nil {
		log.Printf("Ошибка при получении ответа: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось получить ответы"})
		return
	}
	c.JSON(http.StatusOK, answers)

}
func UpdateAnswer(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionId, err := extractQuestionID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	answerId, err := extractAnswerID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	err = service.AnswerChek(creatorId, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при проверке ответа на существование: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Ответ не найден"})
		return
	}

	var updateAnswer models.AnswerRequest

	err = c.ShouldBindJSON(&updateAnswer)

	if err != nil {
		log.Printf("Ошибка при извлечении данных ответа и переводе их в json: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	err = service.AnswerUpdate(updateAnswer, creatorId, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при обновлении данных ответа: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось обновить ответ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Ответ успешно обновлен"})

}
func DeleteAnswer(c *gin.Context) {
	formId, err := extractFormID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id формы %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionId, err := extractQuestionID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	answerId, err := extractAnswerID(c)
	if err != nil {
		log.Printf("Ошибка при извлечении id вопроса %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return
	}

	creatorId, err := extractUserID(c)
	if err != nil {
		log.Printf("Не удалось извлечь id пользователя: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	err = service.AnswerChek(creatorId, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при проверке вопроса %v", err)
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Ответ не найден"})
		return
	}

	_, err = service.AnswerDelete(creatorId, formId, questionId, answerId)
	if err != nil {
		log.Printf("Ошибка при удалении вопроса %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить ответ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Ответ успешно удален"})

}
