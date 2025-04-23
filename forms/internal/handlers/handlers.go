package handlers

import (
	"database/sql"
	"forms/internal/models"
	"forms/internal/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateForm(c *gin.Context) {
	var form models.FormRequest

	err := c.ShouldBindJSON(&form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	var formId int
	formId, link, err := service.FormCreate(form, creatorId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось создать форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"form_id": formId, "link": link})
}

func GetForm(c *gin.Context) {

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	formIdstr := c.Param("id")

	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}
	var form models.Form

	creatorId, ok := creatorIdfl.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный тип id пользователя"})
		return
	}

	form, err = service.FormGet(creatorId, formId)
	// Проверка на наличие возвращаемого значения, если форм нет - сработает условие цикла
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка при получении формы"})
		return
	}

	c.JSON(http.StatusOK, form)
}

func GetForms(c *gin.Context) {
	creatorIdfl, ok := c.Get("id")

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный тип id пользователя"})
		return
	}

	forms, err := service.FormsGet(creatorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось найти анкеты"})
		return
	}
	c.JSON(http.StatusOK, forms)
}

func UpdateForm(c *gin.Context) {
	formIdstr := c.Param("id")

	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	var updateForm models.FormRequest

	err = c.ShouldBindJSON(&updateForm)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.FormChek(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	err = service.FormUpdate(updateForm, creatorId, formId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Ошибка обновления данных"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно обновлена"})
}

func DeleteForm(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorIdfl, ok := c.Get("id")

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}

	creatorId, ok := creatorIdfl.(int)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	// Проверка на существование формы для удаления, нужен id пользователя и id формы
	err = service.FormChek(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	// Удаление формы, для этого требуется id формы и пользователя из-за нужды нахождения формы для удаления с помощью sql-запроса
	_, err = service.FormDelete(formId, creatorId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить форму"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Форма успешно удалена"})
}

func CreateQuestion(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	var question models.QuestionRequest

	err = c.ShouldBindJSON(&question)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	err = service.FormChek(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	questionId, err := service.QuestionCreate(question, formId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось создать вопрос"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Вопрос успешно создан",
		"id вопроса": questionId})

}

func GetQuestions(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.FormChek(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Форма не найдена"})
		return
	}

	questions, err := service.QuestionsGet(creatorId, formId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось получить вопросы"})
		return
	}
	c.JSON(http.StatusOK, questions)

}

func GetQuestion(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.QuestionChek(creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Вопрос не найден"})
		return
	}

	var question models.Question

	question, err = service.QuestionGet(creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось получить вопрос"})
		return
	}
	c.JSON(http.StatusOK, question)

}

func UpdateQuestion(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
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
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	err = service.QuestionUpdate(updateQuestion, creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось обновить вопрос"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Вопрос успешно обновлен"})

}

func DeleteQuestion(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.QuestionChek(creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Вопрос не найден"})
		return
	}

	_, err = service.QuestionDelete(creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить вопрос"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Вопрос успешно удален"})

}

func CreateAnswer(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
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
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	answerId, err := service.AnswerCreate(answer, formId, questionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось создать ответ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Ответ успешно создан",
		"id ответа": answerId})

}

func GetAnswer(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	answerIdstr := c.Param("answer_id")
	// Конвертируем в численный тип данных строку с id для проверки
	answerId, err := strconv.Atoi(answerIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.AnswerChek(creatorId, formId, questionId, answerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Ответ не найден"})
		return
	}

	var answer models.Answer

	answer, err = service.AnswerGet(creatorId, formId, questionId, answerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось получить ответ"})
		return
	}
	c.JSON(http.StatusOK, answer)

}

func GetAnswers(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	answerIdstr := c.Param("answer_id")
	answerId, err := strconv.Atoi(answerIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return
	}

	err = service.AnswerChek(creatorId, formId, questionId, answerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Вопрос не найден"})
		return
	}

	var answers []models.Answer

	answers, err = service.AnswersGet(creatorId, formId, questionId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось получить ответы"})
		return
	}
	c.JSON(http.StatusOK, answers)

}
func UpdateAnswer(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	answerIdstr := c.Param("answer_id")
	// Конвертируем в численный тип данных строку с id для проверки
	answerId, err := strconv.Atoi(answerIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.AnswerChek(creatorId, formId, questionId, answerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Ответ не найден"})
		return
	}

	var updateAnswer models.AnswerRequest

	err = c.ShouldBindJSON(&updateAnswer)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неправильный формат данных"})
		return
	}

	err = service.AnswerUpdate(updateAnswer, creatorId, formId, questionId, answerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось обновить ответ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Ответ успешно обновлен"})

}
func DeleteAnswer(c *gin.Context) {
	formIdstr := c.Param("id")
	// Конвертируем в численный тип данных строку с id для проверки
	formId, err := strconv.Atoi(formIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id формы"})
		return
	}

	questionIdstr := c.Param("question_id")
	// Конвертируем в численный тип данных строку с id для проверки
	questionId, err := strconv.Atoi(questionIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id вопроса"})
		return
	}

	answerIdstr := c.Param("answer_id")
	// Конвертируем в численный тип данных строку с id для проверки
	answerId, err := strconv.Atoi(answerIdstr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "Неверный id ответа"})
		return
	}

	creatorIdfl, ok := c.Get("id")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "id пользователя не найден"})
		return
	}
	creatorId, ok := creatorIdfl.(int)

	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"Ошибка": "неверный тип id"})
		return
	}

	err = service.AnswerChek(creatorId, formId, questionId, answerId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"Ошибка": "Ответ не найден"})
		return
	}

	_, err = service.AnswerDelete(creatorId, formId, questionId, answerId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Ошибка": "Не удалось удалить ответ"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Сообщение": "Ответ успешно удален"})

}
