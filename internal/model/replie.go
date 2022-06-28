package model

type FullReplie struct {
	ID           int    `json:"id" db:"id"`
	UserID       int    `json:"userId" db:"userId"`
	MessageID    int    `json:"messageId" db:"messageId"`
	Title        string `json:"title" db:"title"`
	UserFullname string `json:"userFullname" db:"fullname"`
	UserImageURL string `json:"userImageUrl" db:"imageurl"`
}
