package model

type FullMessage struct {
	ID              int    `json:"id" db:"messageId"`
	Title           string `json:"title" db:"messageTitle"`
	MessageURL      string `json:"messageUrl" db:"messageUrl"`
	MessageImageURL string `json:"messageImageUrl" db:"messageImageUrl"`

	ChannelName     string `json:"channelName" db:"channelName"`
	ChannelTitle    string `json:"channelTitle" db:"channelTitle"`
	ChannelImageURL string `json:"channelImageUrl" db:"channelImageUrl"`

	UserID       int    `json:"userId" db:"userId"`
	UserFullname string `json:"userFullname" db:"userFullname"`
	UserImageURL string `json:"userImageUrl" db:"userImageUrl"`

	RepliesCount int          `json:"repliesCount" db:"count"`
	Replies      []FullReplie `json:"replies"`
}
