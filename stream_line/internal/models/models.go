package models

import "time"

type FormResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
}

type FormFromDB struct {
	ID          string
	Title       string
	Description string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	CreatorID   string
}
