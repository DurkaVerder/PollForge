package models

type Answer struct {
	Title string `json:"title"`
	Count uint64 `json:"count"`
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
	ID         int
	QuestionID int    `json:"question_id"`
	Title      string `json:"title"`
	Count      uint64 `json:"count"`
}
