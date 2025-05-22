package models

import "time"

type FormRequest struct {
	ThemeId	int       `json:"theme_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	PrivateKey  bool      `json:"private_key"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
}


type QuestionRequest struct {
	Title       string `json:"title" binding:"required"`
	NumberOrder int    `json:"number_order" binding:"required"`
	Required    bool   `json:"required"`
}


type AnswerRequest struct {
	Title       string `json:"title" binding:"required"`
	NumberOrder int    `json:"number_order" binding:"required"`
	Count       int    `json:"count"`
	AnswerId    int    `json:"answer_id"`
}

type Form struct {
	Id          int               `json:"id"`
	CreatorId   int               `json:"user_id"`
	ThemeName string            `json:"theme"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Link        string            `json:"link"`
	PrivateKey  bool              `json:"private_key"`
	ExpiresAt   time.Time         `json:"expires_at"`
	CreatedAt   time.Time         `json:"created_at"`
	Questions   []QuestionOutput  `json:"questions"` 
}

type QuestionOutput struct {
	Id          int      `json:"id"`
	NumberOrder int      `json:"number_order"`
	Title       string   `json:"title"`
	Required    bool     `json:"required"`
	Answers     []Answer `json:"answers"`
}

type Answer struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	NumberOrder int    `json:"number_order"`
	Count       int    `json:"count"`
}