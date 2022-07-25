package model

// @Description Saved message model
type Saved struct {
	ID        int `json:"id" db:"id"`                // Saved id example: 1
	UserID    int `json:"userId" db:"user_id"`       // Saved user id example: 1
	MessageID int `json:"messageId" db:"message_id"` // Saved message id example: 1
}
