package models

import "time"

type StreamLineResponse struct {
	Polls []Polls `json:"polls"`
}

type Polls struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Description   string     `json:"description"`
	Link          string     `json:"link"`
	Likes         Like       `json:"likes"`
	CountComments int        `json:"count_comments"`
	Questions     []Question `json:"questions"`
	CreatorID     int        `json:"creator_id"`
	CreatedAt     string     `json:"created_at"`
	ExpiresAt     string     `json:"expires_at"`
}

type Question struct {
	ID              int      `json:"id"`
	Title           string   `json:"title"`
	TotalCountVotes int      `json:"total_count_votes"`
	Answers         []Answer `json:"answers"`
}

type Answer struct {
	ID         int     `json:"id"`
	Title      string  `json:"title"`
	Percent    float64 `json:"percent"`
	CountVotes int     `json:"count_votes"`
	IsSelected bool    `json:"is_selected"`
}

type Like struct {
	Count   int  `json:"count"`
	IsLiked bool `json:"is_liked"`
}

type LikeFromDB struct {
	Count   int
	IsLiked bool
}
type FormFromDB struct {
	ID            int
	Title         string
	Description   string
	Link          string
	Like          LikeFromDB
	CountComments int
	CreatorID     int
	CreatedAt     time.Time
	ExpiresAt     time.Time
}

type QuestionFromDB struct {
	ID              int
	Title           string
	FormID          int
	NumberOrder     int
	TotalCountVotes int // TotalCountVotes. Not from DB
}

type AnswerFromDB struct {
	ID          int
	Title       string
	QuestionID  int
	NumberOrder int
	CountVotes  int
	IsSelected  bool
	Percent     float64
}
