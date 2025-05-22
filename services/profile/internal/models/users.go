package models

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type UserRequest struct {
	Username string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserProfile struct {
	ID       int    `json:"id"`
	Username string `json:"name"`
	Email    string `json:"email"`
	Bio 	string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type Form struct {
	Id          int       `json:"id"`
	CreatorId   int       `json:"creator_id"`
	ThemeId     int       `json:"theme_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link      string    `json:"link"`
	PrivateKey  bool      `json:"private_key"`
	ExpiresAt   time.Time `json:"expires_at"`
}
