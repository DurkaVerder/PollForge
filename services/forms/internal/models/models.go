package models

import "time"

type FormRequest struct {
	ThemeId	int       `json:"theme_id" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Confidential bool	  `json:"confidential"`
	PrivateKey  bool      `json:"private_key"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
}


type QuestionRequest struct {
	Title       string `json:"title" binding:"required"`
	NumberOrder int    `json:"number_order" binding:"required"`
	Required    bool   `json:"required"`
	MultipleChoice bool `json:"multiple_choice"`
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
	Confidential bool              `json:"confidential"`
	Questions   []QuestionOutput  `json:"questions"` 
}

type QuestionOutput struct {
	Id          int      `json:"id"`
	NumberOrder int      `json:"number_order"`
	Title       string   `json:"title"`
	Required    bool     `json:"required"`
	MultipleChoice bool `json:"multiple_choice"`
	Answers     []Answer `json:"answers"`
}

type Answer struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	NumberOrder int    `json:"number_order"`
	Count       int    `json:"count"`
}

type Theme struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}