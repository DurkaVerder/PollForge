package models

import "time"

type Form struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Likes       Like   `json:"likes"`
	CreatedAt   string `json:"created_at"`
	ExpiresAt   string `json:"expires_at"`
}

type FormResponse struct {
	Forms []Form `json:"forms"`
}

type FormFromDB struct {
	ID          int
	Title       string
	Description string
	Like        LikeFromDB
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

type Like struct {
	Count   int  `json:"count"`
	IsLiked bool `json:"is_liked"`
}

type LikeFromDB struct {
	Count   int
	IsLiked bool
}
