package model

// @Description Full replie model includes all info about replie
type FullReplie struct {
	ID           int    `json:"id" db:"id"`                 // Replies id example: 1
	UserID       int    `json:"userId" db:"userid"`         // Replie user id example: 1
	MessageID    int    `json:"messageId" db:"messageid"`   // Replie message id example: 1
	Title        string `json:"title" db:"title"`           // Replie title example: Yes
	UserFullname string `json:"userFullname" db:"fullname"` // Replie user fullname example: Ivan Petrovich
	UserImageURL string `json:"userImageUrl" db:"imageurl"` // Replie user image url from firebase
}
