package storage

import (
	"database/sql"
	"log"
	"stats/internal/models"

	_ "github.com/lib/pq"
)

const (
	QueryGetQuestions = "SELECT id, title FROM questions WHERE form_id = $1 ORDER BY number_order"
	QueryGetAnswers   = "SELECT a.id, a.title, a.question_id, a.count FROM answers a JOIN questions q ON a.question_id = q.id WHERE q.form_id = $1 ORDER BY q.number_order, a.number_order"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) GetQuestions(formID string) ([]models.QuestionFromDB, error) {
	rows, err := p.db.Query(QueryGetQuestions, formID)
	if err != nil {
		log.Printf("GetQuestions: Ошибка при выполнении запроса: %v\n", err)
		return nil, err
	}

	defer rows.Close()
	var questions []models.QuestionFromDB
	for rows.Next() {
		var question models.QuestionFromDB
		if err := rows.Scan(&question.ID, &question.Title); err != nil {
			log.Printf("GetQuestions: Ошибка при сканировании строки: %v\n", err)
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (p *Postgres) GetAnswers(formID string) ([]models.AnswerFromDB, error) {
	rows, err := p.db.Query(QueryGetAnswers, formID)
	if err != nil {
		log.Printf("GetAnswers: Ошибка при выполнении запроса: %v\n", err)
		return nil, err
	}

	defer rows.Close()
	var answers []models.AnswerFromDB
	for rows.Next() {
		var answer models.AnswerFromDB
		if err := rows.Scan(&answer.ID, &answer.Title, &answer.QuestionID, &answer.Count); err != nil {
			log.Printf("GetAnswers: Ошибка при сканировании строки: %v\n", err)
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func (p *Postgres) Close() {
	if err := p.db.Close(); err != nil {
		log.Printf("Close: Ошибка при закрытии соединения с БД: %v\n", err)
	}
}
