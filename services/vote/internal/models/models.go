package models

type VoteRequest struct {
	ID       int  `json:"id"`
	IsUpVote bool `json:"is_up_vote"`
}

type Vote struct {
	ID       int
	IsUpVote bool
	UserID   int
}

type LikeRequest struct {
	ID       int  `json:"id"`
	IsUpLike bool `json:"is_up_like"`
}

type Like struct {
	ID       int
	UserID   int
	IsUpLike bool
}
