package storage

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	maxRetries                  = 3
	QueryUpdateCountAnswerPlus  = "UPDATE answers SET count = count + 1 WHERE id = $1"
	QueryUpdateCountAnswerMinus = "UPDATE answers SET count = count - 1 WHERE id = $1 AND count > 0"

	QueryInsertAnswerChoice = "INSERT INTO answers_chosen (answer_id, user_id) VALUES ($1, $2) ON CONFLICT (answer_id, user_id) DO NOTHING"
	QueryDeleteAnswerChoice = "DELETE FROM answers_chosen WHERE answer_id = $1 AND user_id = $2"
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) UpdateCountAnswerUp(answerId, userID int) error {
	for i := 0; i < maxRetries; i++ {
		tx, err := p.db.Begin()
		if err != nil {
			log.Printf("UpdateCountAnswerUp: Ошибка при начале транзакции: %v\n", err)
			return err
		}

		result, err := tx.Exec(QueryUpdateCountAnswerPlus, answerId)
		if err != nil {
			log.Printf("UpdateCountAnswerUp: Ошибка при выполнении запроса: %v\n", err)

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerUp: Ошибка при откате транзакции: %v\n", err)
			}

			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("UpdateCountAnswerUp: Ошибка при получении количества затронутых строк: %v\n", err)

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerUp: Ошибка при откате транзакции: %v\n", err)
			}

			return err
		}

		if rowsAffected == 0 {
			log.Printf("UpdateCountAnswerUp: Конфликт или некорректные данные, строки не обновлены")

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerUp: Ошибка при откате транзакции: %v\n", err)
			}

			return fmt.Errorf("no rows updated")
		}

		if err := p.createAnswerChoice(tx, answerId, userID); err != nil {
			log.Printf("UpdateCountAnswerUp: Ошибка при создании выбора ответа: %v\n", err)
			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerUp: Ошибка при откате транзакции: %v\n", err)
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("UpdateCountAnswerUp: Ошибка при коммите транзакции: %v\n", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("UpdateCountAnswerUp: failed after %d retries", maxRetries)
}

func (p *Postgres) UpdateCountAnswerDown(answerId, userID int) error {
	for i := 0; i < maxRetries; i++ {
		tx, err := p.db.Begin()
		if err != nil {
			log.Printf("UpdateCountAnswerDown: Ошибка при начале транзакции: %v\n", err)
			return err
		}

		result, err := tx.Exec(QueryUpdateCountAnswerMinus, answerId)
		if err != nil {
			log.Printf("UpdateCountAnswerDown: Ошибка при выполнении запроса: %v\n", err)

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerDown: Ошибка при откате транзакции: %v\n", err)
			}

			continue
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Printf("UpdateCountAnswerDown: Ошибка при получении количества затронутых строк: %v\n", err)

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerDown: Ошибка при откате транзакции: %v\n", err)
			}

			return err
		}

		if rowsAffected == 0 {
			log.Printf("UpdateCountAnswerDown: Конфликт или некорректные данные, строки не обновлены")

			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerDown: Ошибка при откате транзакции: %v\n", err)
			}

			return fmt.Errorf("no rows updated")
		}

		if err := p.deleteAnswerChoice(tx, answerId, userID); err != nil {
			log.Printf("UpdateCountAnswerDown: Ошибка при удалении выбора ответа: %v\n", err)
			if err := tx.Rollback(); err != nil {
				log.Printf("UpdateCountAnswerDown: Ошибка при откате транзакции: %v\n", err)
			}
			continue
		}

		if err := tx.Commit(); err != nil {
			log.Printf("UpdateCountAnswerDown: Ошибка при коммите транзакции: %v\n", err)
			continue
		}

		return nil
	}

	return fmt.Errorf("UpdateCountAnswerDown: failed after %d retries", maxRetries)
}

func (p *Postgres) createAnswerChoice(tx *sql.Tx, answerId int, userId int) error {
	_, err := tx.Exec(QueryInsertAnswerChoice, answerId, userId)
	if err != nil {
		log.Printf("createAnswerChoice: Ошибка при выполнении запроса: %v\n", err)
		return err
	}

	return nil
}

func (p *Postgres) deleteAnswerChoice(tx *sql.Tx, answerId int, userId int) error {
	_, err := tx.Exec(QueryDeleteAnswerChoice, answerId, userId)
	if err != nil {
		log.Printf("deleteAnswerChoice: Ошибка при выполнении запроса: %v\n", err)
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
