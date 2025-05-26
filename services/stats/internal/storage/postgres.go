package storage

import (
	"database/sql"
	"log"
	"stats/internal/models"

	"github.com/lib/pq"
)

const (
	QueryGetQuestions  = "SELECT id, title FROM questions WHERE form_id = $1 ORDER BY number_order"
	QueryGetAnswers    = "SELECT a.id, a.title, a.question_id, a.count FROM answers a JOIN questions q ON a.question_id = q.id  WHERE q.form_id = $1 ORDER BY q.number_order, a.number_order"
	QueryGetTimeChosen = "SELECT answer_id, created_at FROM answers_chosen WHERE answer_id = ANY($1) "

	QueryGetProfileStats = `
SELECT 
    COUNT(DISTINCT f.id) AS count_polls,
    COALESCE(SUM(a.count), 0) AS count_votes
FROM 
    forms f
LEFT JOIN 
    questions q ON f.id = q.form_id
LEFT JOIN 
    answers a ON q.id = a.question_id
WHERE 
    f.creator_id = $1 
    AND f.private_key = FALSE;`

	QueryGetCountComments = `
SELECT
	COUNT(DISTINCT c.id) AS count_comments
FROM
	comments c
JOIN
	forms f ON c.form_id = f.id
WHERE 	
	f.creator_id = $1
	AND f.private_key = FALSE;`

	QueryGetThemeStats = `
SELECT
    t.name,
    t.description,
    COUNT(DISTINCT f.id) AS count_polls,
    COALESCE(SUM(a.count), 0) AS count_votes
FROM
    themes t
JOIN
    forms f ON t.id = f.theme_id
LEFT JOIN
    questions q ON f.id = q.form_id
LEFT JOIN
    answers a ON q.id = a.question_id
WHERE
    f.creator_id = $1
GROUP BY
    t.id
ORDER BY
    t.name;
`
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

func (p *Postgres) GetTimeChosen(answerIDs []int) ([]models.TimeChosenFromDB, error) {
	rows, err := p.db.Query(QueryGetTimeChosen, pq.Array(answerIDs))
	if err != nil {
		log.Printf("GetTimeChosen: Ошибка при выполнении запроса: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var times []models.TimeChosenFromDB
	for rows.Next() {
		var time models.TimeChosenFromDB
		if err := rows.Scan(&time.IdAnswer, &time.Time); err != nil {
			log.Printf("GetTimeChosen: Ошибка при сканировании строки: %v\n", err)
			return nil, err
		}
		times = append(times, time)
	}
	if err := rows.Err(); err != nil {
		log.Printf("GetTimeChosen: Ошибка при итерации по строкам: %v\n", err)
		return nil, err
	}

	return times, nil
}

func (p *Postgres) GetProfileStats(userID string) (models.ProfileStatsFromDB, error) {
	row := p.db.QueryRow(QueryGetProfileStats, userID)
	var profileStats models.ProfileStatsFromDB
	if err := row.Scan(&profileStats.CountCreated, &profileStats.CountAnswered); err != nil {
		if err == sql.ErrNoRows {
			return models.ProfileStatsFromDB{}, nil
		}
		log.Printf("GetProfileStats: Ошибка при выполнении запроса: %v\n", err)
		return models.ProfileStatsFromDB{}, err
	}

	row = p.db.QueryRow(QueryGetCountComments, userID)
	if err := row.Scan(&profileStats.CountComments); err != nil {
		if err == sql.ErrNoRows {
			profileStats.CountComments = 0
		} else {
			log.Printf("GetProfileStats: Ошибка при выполнении запроса: %v\n", err)
			return models.ProfileStatsFromDB{}, err
		}
	}

	return profileStats, nil
}

func (p *Postgres) GetThemeStats(userID string) ([]models.Theme, error) {
	rows, err := p.db.Query(QueryGetThemeStats, userID)
	if err != nil {
		log.Printf("GetThemeStats: Ошибка при выполнении запроса: %v\n", err)
		return nil, err
	}
	defer rows.Close()
	var themes []models.Theme
	for rows.Next() {
		var theme models.Theme
		if err := rows.Scan(&theme.Name, &theme.Description, &theme.CountPolls, &theme.CountVotes); err != nil {
			log.Printf("GetThemeStats: Ошибка при сканировании строки: %v\n", err)
			return nil, err
		}
		themes = append(themes, theme)
	}

	return themes, nil
}
