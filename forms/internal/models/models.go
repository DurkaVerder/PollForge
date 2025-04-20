package models

import "time"

type FormRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description" binding:"required"`
	PrivateKey  bool      `json:"private_key" binding:"required"`
	ExpiresAt   time.Time `json:"expires_at" binding:"required"`
}

type Form struct {
	Id          int       `json:"id"`
	CreatorId   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	PrivateKey  bool      `json:"private_key"`
	ExpiresAt   time.Time `json:"expires_at"`
}
