package models

import "time"

type Message struct {
	ID          uint64
	SenderID    uint64
	RecipientID uint64
	Content     string
	SentAt      time.Time
}

type NewMessageReq struct{
	SenderID  uint64 `json:"sender"`
	RecipentID uint64 `json:"recipent"`
	Message string  `json:"message"`
}