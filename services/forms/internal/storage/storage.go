package storage

import (
	"database/sql"
	"forms/internal/models"
	"log"
	"os"
	"time"

	"github.com/google/uuid"

	_ "github.com/lib/pq"
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
	err := Db.QueryRow(queryChek, formId, creatorId).Scan(&existId)
	if err != nil {
		log.Printf("Ошибка при запросе на проверку наличия формы: %v", err)
		return err
	}
	return err
}

func FormCreateRequest(form models.FormRequest, creatorId int) (int, string, error) {

	link := uuid.New().String()
	query := `INSERT INTO forms (creator_id, theme_id, title, description, link, private_key, expires_at, created_at) 
			  VALUES($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	var formId int
	createdAt := time.Now().Local()
	err := Db.QueryRow(query, creatorId, form.ThemeId, form.Title, form.Description, link, form.PrivateKey, form.ExpiresAt, createdAt).Scan(&formId)
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

func FormUpdateRequest(updateForm models.FormRequest, creatorId int, formId int) error {

	query := "UPDATE forms SET title = $1, description = $2, private_key = $3, expires_at = $4, theme_id = $5 WHERE id = $6 AND creator_id = $7"
	_, err := Db.Exec(query, updateForm.Title, updateForm.Description, updateForm.PrivateKey, updateForm.ExpiresAt, updateForm.ThemeId, formId, creatorId)
	if err != nil {
		log.Printf("Ошибка при запросе обновления формы: %v", err)
		return err
	}
	return err
}

func FormGetRequest(creatorId int, formId int) (models.Form, error) {
	var form models.Form
	query := `
		SELECT f.id, t.name, f.title, f.description, f.link, f.private_key, f.expires_at, f.created_at
		FROM forms AS f LEFT JOIN themes AS t ON f.theme_id = t.id
		WHERE f.id = $1 AND f.creator_id = $2
		`
	err := Db.QueryRow(query, formId, creatorId).Scan(
		&form.Id,
		&form.ThemeName,
		&form.Title,
		&form.Description,
		&form.Link,
		&form.PrivateKey,
		&form.ExpiresAt,
		&form.CreatedAt,
	)
	if err != nil {
		return form, err
	}
	return form, nil
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
func QuestionCreateRequest(question models.QuestionRequest, creatorId int, formId int) (int, error) {
	query := `INSERT INTO questions (form_id, creator_id, title, number_order, required) 
			  VALUES($1, $2, $3, $4, $5) RETURNING id`

	
	var questionId int
	err := Db.QueryRow(query, formId, creatorId, question.Title, question.NumberOrder, question.Required).Scan(&questionId)
	if err != nil {
		log.Printf("Ошибка при запросе создания вопроса: %v", err)
		return questionId, err
	}
	return questionId, err
}
func QuestionDeleteRequest(creator_id int, formId int, questionId int) (sql.Result, error) {

	query := "DELETE FROM questions WHERE id = $1 AND form_id = $2 and creator_id = $3"
	_, err := Db.Exec(query, questionId, formId, creator_id)
	if err != nil {
		log.Printf("Ошибка при запросе удаления вопроса: %v", err)
		return nil, err
	}
	return nil, err
}

func QuestionsGetRequest(creator_id int, formId int) (*sql.Rows, error) {
	query := `
			SELECT
			q.id,
			q.form_id,
			q.title,
			q.number_order,
			q.required,
			COALESCE(a.id, 0)           AS answer_id,
			COALESCE(a.title, '')       AS answer_title,
			COALESCE(a.number_order, 0) AS answer_order,
			COALESCE(a.count, 0)        AS answer_count
			FROM questions AS q
			LEFT JOIN answers AS a
			ON q.id = a.question_id
			WHERE q.creator_id = $1
			AND q.form_id    = $2
			ORDER BY
			q.number_order,
			a.number_order
`

	rows, err := Db.Query(query, creator_id, formId)

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

func AnswerCreateRequest(answer models.AnswerRequest, creatorId int, formId int, questionId int) (int, error) {
	query := `INSERT INTO answers (question_id, form_id, creator_id, title, number_order, count) 
			  VALUES($1, $2, $3, $4, $5, $6) RETURNING id`

	var answerId int
	err := Db.QueryRow(query, questionId, formId, creatorId, answer.Title, answer.NumberOrder, answer.Count).Scan(&answerId)
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

func QuestionsWithAnswersGet(formId, creatorId int) ([]models.QuestionOutput, error) {
	query := `
			SELECT
			q.id, 
			q.number_order,
			q.title,
			q.required,
			a.id,
			a.title,
			a.number_order,
			a.count
			FROM questions q
			LEFT JOIN answers a ON q.id = a.question_id
			WHERE q.form_id = $1 AND q.creator_id = $2
			ORDER BY q.number_order, a.number_order`

	rows, err := Db.Query(query, formId, creatorId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	questionMap := make(map[int]*models.QuestionOutput)
	var orderedIDs []int

	for rows.Next() {
		var (
			qID, qOrder    int
			qTitle         string
			qRequired      bool
			aID            sql.NullInt64
			aTitle         sql.NullString
			aOrder, aCount sql.NullInt64
		)

		err := rows.Scan(&qID, &qOrder, &qTitle, &qRequired, &aID, &aTitle, &aOrder, &aCount)
		if err != nil {
			return nil, err
		}

		q, exists := questionMap[qID]
		if !exists {
			q = &models.QuestionOutput{
				Id:          qID,
				NumberOrder: qOrder,
				Title:       qTitle,
				Required:    qRequired,
				Answers:     []models.Answer{},
			}
			questionMap[qID] = q
			orderedIDs = append(orderedIDs, qID)
		}

		if aID.Valid {
			q.Answers = append(q.Answers, models.Answer{
				Id:          int(aID.Int64),
				Title:       aTitle.String,
				NumberOrder: int(aOrder.Int64),
				Count:       int(aCount.Int64),
			})
		}
	}

	questions := make([]models.QuestionOutput, 0, len(orderedIDs))
	for _, id := range orderedIDs {
		questions = append(questions, *questionMap[id])
	}
	return questions, nil
}

func GetFormByLinkRequest(link string) (models.Form, error) {
	var form models.Form
	query := `
		SELECT f.id, t.name, f.creator_id, f.title, f.description, f.link, f.private_key, f.expires_at, f.created_at
		FROM forms AS f LEFT JOIN themes AS t ON f.theme_id = t.id
		WHERE link = $1
		`
	err := Db.QueryRow(query, link).Scan(
		&form.Id,
		&form.ThemeName,
		&form.CreatorId,
		&form.Title,
		&form.Description,
		&form.Link,
		&form.PrivateKey,
		&form.ExpiresAt,
		&form.CreatedAt,
	)
	if err != nil {
		return form, err
	}
	return form, nil
}


func GetThemesRequest() ([]models.Theme, error) {
	var themes []models.Theme
	query := "SELECT id, name, description FROM themes"
	rows, err := Db.Query(query)
	if err != nil {
		log.Printf("Ошибка при запросе получения тем: %v", err)
		return themes, err
	}
	defer rows.Close()

	for rows.Next() {
		var theme models.Theme
		err := rows.Scan(&theme.Id, &theme.Name, &theme.Description)
		if err != nil {
			log.Printf("Ошибка при сканировании темы: %v", err)
			return themes, err
		}
		themes = append(themes, theme)
	}
	return themes, nil
}