package models

import "time"

type Comment struct {
	Id	   int       `json:"id"`
	FormID   int    `json:"form_id"`
	UserId int    `json:"user_id"`
	Description string `json:"description"`
	CreatedAt  time.Time `json:"created_at"`
}
