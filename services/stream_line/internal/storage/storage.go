package storage

import (
	"database/sql"
	"stream_line/internal/models"
	"time"

	"github.com/lib/pq"
)

const (
	GetOtherFormsQuery = `SELECT f.id, f.title, t.name, f.description, f.creator_id, f.link, f.count_likes, EXISTS(SELECT * FROM likes_forms WHERE user_id = $1 AND form_id = f.id) AS is_liked, (SELECT COUNT(id) FROM comments c WHERE c.form_id = f.id) AS count_votes, f.created_at, f.expires_at FROM forms f LEFT JOIN themes t ON t.id = f.theme_id WHERE f.expires_at > NOW() AND f.private_key = false AND f.created_at < $2 ORDER BY f.created_at DESC LIMIT $3`
	GetQuestionQuery   = `SELECT id, title, form_id, number_order FROM questions WHERE form_id = ANY($1)`
	GetAnswerQuery     = `SELECT a.id, a.title, a.question_id, a.number_order, a.count, EXISTS(SELECT * FROM answers_chosen WHERE user_id = $2 AND answer_id = a.id) AS is_selected FROM answers a WHERE a.question_id = ANY($1)`
	GetFormsQuery      = `SELECT f.id, f.title, t.name, f.description, f.creator_id, f.link, f.count_likes, EXISTS(SELECT * FROM likes_forms WHERE user_id = $1 AND form_id = f.id) AS is_liked, (SELECT COUNT(id) FROM comments c WHERE c.form_id = f.id) AS count_votes, f.created_at, f.expires_at FROM forms f LEFT JOIN themes t ON t.id = f.theme_id WHERE f.expires_at > NOW() AND f.link = $2`
)

type Postgres struct {
	db *sql.DB
}

func NewPostgres(db *sql.DB) *Postgres {
	return &Postgres{
		db: db,
	}
}

func (p *Postgres) GetFormsByOtherUserIDWithCountLikesAndComments(userID string, cursor time.Time, limit int) ([]models.FormFromDB, error) {
	rows, err := p.db.Query(GetOtherFormsQuery, userID, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []models.FormFromDB
	for rows.Next() {
		var form models.FormFromDB
		if err := rows.Scan(&form.ID, &form.Title, &form.Theme, &form.Description, &form.CreatorID, &form.Link, &form.Like.Count, &form.Like.IsLiked, &form.CountComments, &form.CreatedAt, &form.ExpiresAt); err != nil {
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

func (p *Postgres) GetFormsByOtherUserIDWithCountLikesAndCommentsByLink(userID, pollLink string) ([]models.FormFromDB, error) {
	rows, err := p.db.Query(GetFormsQuery, userID, pollLink)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var forms []models.FormFromDB
	for rows.Next() {
		var form models.FormFromDB
		if err := rows.Scan(&form.ID, &form.Title, &form.Theme, &form.Description, &form.CreatorID, &form.Link, &form.Like.Count, &form.Like.IsLiked, &form.CountComments, &form.CreatedAt, &form.ExpiresAt); err != nil {
			return nil, err
		}
		forms = append(forms, form)
	}

	return forms, nil
}
