package storage

import (
	"database/sql"
	"stream_line/internal/models"

	"github.com/lib/pq"
)

const (
	GetOtherFormsQuery = `SELECT f.id, f.title, f.description, f.link, COALESCE(l.count, 0) AS count, EXISTS(SELECT * FROM likes_forms WHERE user_id = $1 AND form_id = f.id) AS is_liked, (SELECT COUNT(id) FROM comments c WHERE c.form_id = f.id) AS count_votes, f.created_at, f.expires_at FROM forms f LEFT JOIN likes l ON l.form_id = f.id WHERE f.expires_at > NOW() AND f.private_key = false LIMIT 10`
	GetQuestionQuery   = `SELECT id, title, form_id, number_order FROM questions WHERE form_id IN ANY($1)`
	GetAnswerQuery     = `SELECT a.id, a.title, a.question_id, a.number_order, a.count, EXISTS(SELECT * FROM answers_chosen WHERE user_id = $2 AND answer_id = a.id) AS is_selected FROM answers a WHERE a.question_id IN ANY($1)`
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) GetFormsByOtherUserIDWithCountLikesAndComments(userID string) ([]models.FormFromDB, error) {
	rows, err := p.db.Query(GetOtherFormsQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []models.FormFromDB
	for rows.Next() {
		var form models.FormFromDB
		if err := rows.Scan(&form.ID, &form.Title, &form.Description, &form.Link, &form.Like.Count, &form.Like.IsLiked, &form.CountVotes, &form.CreatedAt, &form.ExpiresAt); err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	return forms, nil
}

func (p *Postgres) GetQuestionsByFormsID(formIDs []int) ([]models.QuestionFromDB, error) {
	rows, err := p.db.Query(GetQuestionQuery, pq.Array(formIDs))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []models.QuestionFromDB
	for rows.Next() {
		var question models.QuestionFromDB
		if err := rows.Scan(&question.ID, &question.Title, &question.FormID, &question.NumberOrder); err != nil {
			return nil, err
		}
		questions = append(questions, question)
	}

	return questions, nil
}

func (p *Postgres) GetAnswersByQuestionsID(questionIDs []int, userID string) ([]models.AnswerFromDB, error) {
	rows, err := p.db.Query(GetAnswerQuery, pq.Array(questionIDs), userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var answers []models.AnswerFromDB
	for rows.Next() {
		var answer models.AnswerFromDB
		if err := rows.Scan(&answer.ID, &answer.Title, &answer.QuestionID, &answer.NumberOrder, &answer.CountVotes, &answer.IsSelected); err != nil {
			return nil, err
		}
		answers = append(answers, answer)
	}

	return answers, nil
}
