package model

type Saved struct {
	ID        int `json:"id" db:"id"`
	UserID    int `json:"userId" db:"user_id"`
	MessageID int `json:"messageId" db:"message_id"`
}
