package models

import "time"

type Answer struct {
	Title        string   `json:"title"`
	Count        uint64   `json:"count"`
	TimeSelected []string `json:"time_selected"`
}

type PollStats struct {
	QuestionTitle string   `json:"question_title"`
	Answers       []Answer `json:"answers"`
}

type PollStatsResponse struct {
	Stats []PollStats `json:"stats"`
}

type QuestionFromDB struct {
	ID    int
	Title string `json:"title"`
}

type AnswerFromDB struct {
	ID           int
	QuestionID   int
	Title        string
	Count        uint64
	TimeSelected []string
}

type TimeChosenFromDB struct {
	IdAnswer int
	Time     time.Time
}

type ProfileStatsRequest struct {
	CountCreatedPolls    int     `json:"count_created_polls"`
	CountAnsweredPolls   int     `json:"count_answered_polls"`
	CountCommentsByPolls int     `json:"count_comments_by_poll"`
	Themes               []Theme `json:"themes"`
}

type Theme struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	CountPolls  int    `json:"count_polls"`
	CountVotes  int    `json:"count_votes"`
}

type ProfileStatsFromDB struct {
	CountCreated  int
	CountAnswered int
	CountComments int
}
