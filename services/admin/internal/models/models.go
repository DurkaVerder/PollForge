package models

import "time"

type UserRequest struct {
	ID       int    `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	ID        int    `json:"id"`
	Username  string `json:"name"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
	IsBanned  bool   `json:"is_banned"`
}

type ToogleBan struct {
	IsBanned bool `json:"is_banned"`
}

type Form struct {
	Id          int              `json:"id"`
	CreatorId   int              `json:"user_id"`
	Title       string           `json:"title"`
	Description string           `json:"description"`
	Link        string           `json:"link"`
	PrivateKey  bool             `json:"private_key"`
	ExpiresAt   time.Time        `json:"expires_at"`
	Questions   []QuestionOutput `json:"questions"`
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