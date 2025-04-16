package storage

import (
	"database/sql"
	"log"
	"question/models"
)

const (
	QueryGetQuestions      = "SELECT id, title FROM questions WHERE form_id = $1 ORDER BY number_order"
	QueryGetAnswers        = "SELECT a.id, a.title, a.question_id FROM answers a JOIN questions q ON a.question_id = q.id WHERE q.form_id = $1 ORDER BY q.number_order, a.number_order"
	QueryUpdateCountAnswer = "UPDATE answers SET count = count + 1 WHERE id IN ($1)"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) GetQuestions(formID int) ([]models.QuestionFromDB, error) {
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

func (p *Postgres) GetAnswers(formID int) ([]models.AnswerFromDB, error) {
	rows, err := p.db.Query(QueryGetAnswers, formID)
	if err != nil {
		log.Printf("GetAnswers: Ошибка при выполнении запроса: %v\n", err)
		return nil, err
	}

	defer rows.Close()
	var answers []models.AnswerFromDB
	for rows.Next() {
		var answer models.AnswerFromDB
		if err := rows.Scan(&answer.ID, &answer.Title, &answer.QuestionID); err != nil {
			log.Printf("GetAnswers: Ошибка при сканировании строки: %v\n", err)
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}

func (p *Postgres) UpdateCountAnswer(ids string) error {
	tx, err := p.db.Begin()
	if err != nil {
		log.Printf("UpdateCountAnswer: Ошибка при начале транзакции: %v\n", err)
		return err
	}

	_, err = p.db.Exec(QueryUpdateCountAnswer, ids)
	if err != nil {
		log.Printf("UpdateCountAnswer: Ошибка при выполнении запроса: %v\n", err)
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
