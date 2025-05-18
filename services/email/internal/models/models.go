package models

type MessageKafka struct {
	EventType string `json:"event_type"`
	UserID    string `json:"user_id"`
	Token 	  string `json:"token,omitempty"`
}

