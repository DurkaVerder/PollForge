package storage

import (
	"database/sql"
	"fmt"
	"log"
	"question/models"

	"github.com/lib/pq"
)

const (
	maxRetries               = 3
	QueryGetQuestions        = "SELECT id, title FROM questions WHERE form_id = $1 ORDER BY number_order"
	QueryGetAnswers          = "SELECT a.id, a.title, a.question_id FROM answers a JOIN questions q ON a.question_id = q.id WHERE q.form_id = $1 ORDER BY q.number_order, a.number_order"
	QueryUpdateCountAnswer   = "UPDATE answers SET count = count + 1 WHERE id = ANY($1)"
	QueryExistsUserAnswer    = "SELECT EXISTS(SELECT 1 FROM answered_polls WHERE form_id = $1 AND user_id = $2)"
	QueryInsertAnsweredPolls = "INSERT INTO answered_polls (form_id, user_id) VALUES ($1, $2)"
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

func (p *Postgres) UpdateCountAnswer(ids []int) error {

	for i := 0; i < maxRetries; i++ {
		tx, err := p.db.Begin()
		if err != nil {
			log.Printf("UpdateCountAnswer: Ошибка при начале транзакции: %v\n", err)
			return err
		}

		result, err := tx.Exec(QueryUpdateCountAnswer, pq.Array(ids))
		if err != nil {
			log.Printf("UpdateCountAnswer: Ошибка при выполнении запроса: %v\n", err)

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswer: Ошибка при откате транзакции: %v\n", err)
			}

			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("UpdateCountAnswer: Ошибка при получении количества затронутых строк: %v\n", err)

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswer: Ошибка при откате транзакции: %v\n", err)
			}

			return err
		}

		if rowsAffected == 0 {
			log.Printf("UpdateCountAnswer: Конфликт или некорректные данные, строки не обновлены")

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswer: Ошибка при откате транзакции: %v\n", err)
			}

			return fmt.Errorf("no rows updated")
		}

		if err := tx.Commit(); err != nil {
			log.Printf("UpdateCountAnswer: Ошибка при коммите транзакции: %v\n", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("UpdateCountAnswer: failed after %d retries", maxRetries)
}

func (p *Postgres) ExistsUserAnswer(formId, userId int) (bool, error) {
	row := p.db.QueryRow(QueryExistsUserAnswer, formId, userId)

	exists := false
	if err := row.Scan(&exists); err != nil {
		log.Printf("ExistsUserAnswer: Ошибка при сканировании строки: %v\n", err)
		return exists, err
	}

	return exists, nil
}

func (p *Postgres) CreateAnsweredPolls(formId, userId int) error {
	if _, err := p.db.Exec(QueryInsertAnsweredPolls, formId, userId); err != nil {
		log.Printf("CreateAnsweredPolls: Ошибка при выполнении запроса: %v\n", err)
		return err
	}

	return nil
}

func (p *Postgres) Close() error {
	if err := p.db.Close(); err != nil {
		return err
	}

	return nil
}
