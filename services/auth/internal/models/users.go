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
}

type Form struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
}

type MessageKafka struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Token     string `json:"token,omitempty"`
}

type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type PasswordResetConfirm struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type PasswordReset struct {
	ID        int       `json:"id"`        
	UserID    int       `json:"user_id"`   
	Token     string    `json:"token"`      
	ExpiresAt time.Time `json:"expires_at"` 
}

type RoleAndBan struct {
	Role     string `json:"role"`
	IsBanned bool   `json:"is_banned"`
}