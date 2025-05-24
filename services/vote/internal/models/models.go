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
