package model

// @Description Full message model includes all info about message
type FullMessage struct {
	ID              int    `json:"id" db:"messageid"`                    // Message id example: 1
	Title           string `json:"title" db:"messagetitle"`              // Message title example: Hello, anyone works with Golang?
	MessageURL      string `json:"messageUrl" db:"messageurl"`           // Message url from telegram
	MessageImageURL string `json:"messageImageUrl" db:"messageimageurl"` // Message image url from firebase

	ChannelName     string `json:"channelName" db:"channelname"`         // Channel name example: go_go
	ChannelTitle    string `json:"channelTitle" db:"channeltitle"`       // Channel title example: GO ukrainian community
	ChannelImageURL string `json:"channelImageUrl" db:"channelimageurl"` // Channel image url from firebase

	UserID       int    `json:"userId" db:"userid"`             // User id example: 1
	UserFullname string `json:"userFullname" db:"userfullname"` // User fullname example: Ivan Petrovich
	UserImageURL string `json:"userImageUrl" db:"userimageurl"` // User image url from firebase

	RepliesCount int          `json:"repliesCount" db:"count"` // Replies count example: 50
	Replies      []FullReplie `json:"replies"`                 // Replies
}

type MessageDTO struct {
	ChannelID  int    `db:"channel_id"`
	UserID     int    `db:"user_id"`
	Title      string `db:"title"`
	MessageURL string `db:"message_url"`
	ImageURL   string `db:"imageurl"`
}
