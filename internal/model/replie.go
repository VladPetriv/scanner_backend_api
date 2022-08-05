package model

// @Description Full replie model includes all info about replie
type FullReplie struct {
	ID           int    `json:"id" db:"id"`                     // Replies id example: 1
	MessageID    int    `json:"messageId" db:"message_id"`      // Replie message id example: 1
	Title        string `json:"title" db:"title"`               // Replie title example: Yes
	ImageURL     string `json:"imageurl" db:"imageurl"`         // Replie image url from firebase
	UserID       int    `json:"userId" db:"userid"`             // Replie user id example: 1
	UserFullname string `json:"userFullname" db:"fullname"`     // Replie user fullname example: Ivan Petrovich
	UserImageURL string `json:"userImageUrl" db:"userimageurl"` // Replie user image url from firebase
}

type ReplieDTO struct {
	MessageID int    `db:"message_id"`
	UserID    int    `db:"user_id"`
	Title     string `db:"title"`
	ImageURL  string `db:"imageurl"`
}
