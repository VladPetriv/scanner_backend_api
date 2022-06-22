package model

type Message struct {
	ID        int    `json:"id"`
	UserID    int    `json:"userId"`
	ChannelID int    `json:"channelId"`
	Title     string `json:"title"`
	ImageURL  string `json:"imageUrl"`
}
