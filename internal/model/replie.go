package model

type Replie struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	MessageID int    `json:"messageId"`
	Title     string `json:"title"`
}