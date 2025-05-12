package models

import "time"

type FormRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	PrivateKey  bool      `json:"private_key"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
}

type Form struct {
	Id          int       `json:"id"`
	CreatorId   int       `json:"user_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PrivateKey  bool      `json:"private_key"`
	ExpiresAt   time.Time `json:"expires_at"`
}

type Question struct {
	Id                int    `json:"id"`
	FormId            int    `json:"form_id"`
	NumberOrder       int    `json:"number_order"`
	Title             string `json:"title"`
	Required          bool   `json:"required"`
	AnswerTitle       string `json:"answer_title"`
	AnswerNumberOrder int    `json:"answer_order"`
	AnswerCount       int    `json:"answer_count"`
	AnswerChosen	  bool   `json:"answer_chosen"`
}

type QuestionRequest struct {
	Title       string `json:"title" binding:"required"`
	NumberOrder int    `json:"number_order" binding:"required"`
	Required    bool   `json:"required" binding:"required"`
}

type Answer struct {
	Id          int    `json:"id"`
	QuestionId  int    `json:"question_id"`
	Title       string `json:"title"`
	NumberOrder int    `json:"number_order"`
	Count       int    `json:"count"`
	Chosen      bool   `json:"chosen"`

}

type AnswerRequest struct {
	Title       string `json:"title" binding:"required"`
	NumberOrder int    `json:"number_order" binding:"required"`
	Count       int    `json:"count"`
	AnswerId    int    `json:"answer_id"`
}
