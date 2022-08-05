package model

// @Description Channel model
type Channel struct {
	ID       int    `json:"id"`       // channel id example: 1
	Name     string `json:"name"`     // channel name example: go_go
	Title    string `json:"title"`    // channel title example: GO ukrainian community
	ImageURL string `json:"imageUrl"` // channel image url from firebase
}

type ChannelDTO struct {
	Name     string `json:"Username" db:"name"`
	Title    string `json:"Title" db:"title"`
	ImageURL string `json:"ImageURL" db:"imageurl"`
}
