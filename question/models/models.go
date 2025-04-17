package models


type SubmitAnswer struct {
	ID       int  `json:"id"`
	Selected bool `json:"selected"`
}

type SubmitAnswerRequest struct {
	Answers []SubmitAnswer `json:"answers"`
}

type Answer struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

type Question struct {
	ID      int      `json:"id"`
	Title   string   `json:"title"`
	Answers []Answer `json:"answers"`
}

type QuestionResponse struct {
	Question []Question `json:"questions"`
}

type QuestionFromDB struct {
	ID    int
	Title string
}

type AnswerFromDB struct {
	ID         int
	QuestionID int
	Title      string
}
